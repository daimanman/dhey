/**

 * daixiaoman  获取表单参数
 */
yihg_utils_fun = {
    getParams : function(paramDivId) {
        var params = {};
        var _$ = jQuery || $;
        var _jqueryObj ;
        if((typeof paramDivId ) == "string"){
            _jqueryObj = _$("#" + paramDivId);
        }else if (paramDivId instanceof jQuery){
            _jqueryObj = paramDivId;
        }
        _jqueryObj.find(" input , select ,textarea").each(
            function(i,o) {
                var jqObj = _$(o);
                var type = jqObj.attr("type");
                var val = _$.trim(jqObj.val());
                var name = jqObj.attr("name");
                if(!name){
                    return true;
                }
                if (type == "text" || type == "hidden"
                    || jqObj.is("select") || type=="password") {
                    params[name] = val;
                }

                if (type == "radio") {
                    jqObj.attr("checked") && (params[name] = val);
                }

                if ((type == "checkbox")) {
                    var checkAttr = jqObj.attr("checked");

                    if (checkAttr && !params[name]) {
                        params[name] = [];
                    }
                    checkAttr && (params[name]).push(val);
                }

                if (jqObj.is("textarea")) {
                    params[name] = jqObj.val();
                }
            });
        return params;
    },
    getChildParams:function(paramDivId){
        //获取子层级input select textarea 参数 不包括所有后代
        var params = {};
        var _$ = jQuery || $;
        var _jqueryObj ;
        if((typeof paramDivId ) == "string"){
            _jqueryObj = _$("#" + paramDivId);
        }else if (paramDivId instanceof jQuery){
            _jqueryObj = paramDivId;
        }
        _jqueryObj.find(" input , select ,textarea").each(
            function(i,o) {
                var jqObj = _$(o);
                var deptLength = yihg_utils_fun._getDepth(_jqueryObj,jqObj);
                if(deptLength > 1){
                    console.info(deptLength);
                    return true;
                }
                var type = jqObj.attr("type");
                var val = _$.trim(jqObj.val());
                var name = jqObj.attr("name");
                if(!name){
                    return true;
                }
                if (type == "text" || type == "hidden"
                    || jqObj.is("select") || type=="password") {
                    params[name] = val;
                }

                if (type == "radio") {
                    jqObj.attr("checked") && (params[name] = val);
                }

                if ((type == "checkbox")) {
                    var checkAttr = jqObj.attr("checked");

                    if (checkAttr && !params[name]) {
                        params[name] = [];
                    }
                    checkAttr && (params[name]).push(val);
                }

                if (jqObj.is("textarea")) {
                    params[name] = jqObj.val();
                }
            });
        return params;
    },
    _getDepth:function(curJq,target){
        //获取路径深度
        var length = 0;
        if(!curJq.hasClass("params-group")){
            length = 1;
        }

        var parent =target.parent();
        if(parent.hasClass("params-group")){
            length++;
        }

        while((!parent.is(curJq))){
            parent = parent.parent();
            if(parent.hasClass("params-group")){
                length++;
            }
        }
        return length;

    },
    getGroupParams:function(paramDivId){
        var _$ = jQuery || $;
        var _jqueryObj ;
        if((typeof paramDivId ) == "string"){
            _jqueryObj = _$("#" + paramDivId);
        }else if (paramDivId instanceof jQuery){
            _jqueryObj = paramDivId;
        }


        //获取第一层参数
        var params = yihg_utils_fun.getChildParams(_jqueryObj);

        _jqueryObj.children(".params-group").each(function(i,groupObj){
            var groupJq = $(groupObj);
            var groupName = groupJq.attr("group-name");
            if(groupName){
                if(!params[groupName]){
                    params[groupName] = [];
                }
                (params[groupName]).push(yihg_utils_fun.getGroupParams(groupJq));
            }
        });
        return params;


    },
    resetParams : function(paramDivId, excludesMap) {
        var excludes = excludesMap || {};
        var _$ = jQuery || $;
        _$("#" + paramDivId).find(" input , select,textarea ").each(
            function(i) {
                var jqObj = _$(this);
                var type = jqObj.attr("type");
                var name = jqObj.attr("name");
                if(type == "radio" ||  type == "checkbox"){

                }else{
                    !excludes[name] && jqObj.val("");
                }
            });
    },
    showParams : function(paramDivId, params) {
        var _$ = jQuery || $;
        _$("#" + paramDivId).find(" input , select,textarea ").each(
            function(i) {
                var jqObj = _$(this);
                var name = jqObj.attr("name");
                var type = jqObj.attr("type");
                if(type == "radio" || type == "checkbox"){
                    var inputVal = jqObj.val();
                    if((params[name] != undefined) && ($.trim(inputVal) == (params[name]+"") ) ){
                        jqObj.prop("checked",true);
                    }
                }else{
                    (params[name] != undefined) && jqObj.val(params[name]);
                }
            });
    },
    setDefault:function(paramDivId,params){
        var _$ = jQuery || $;
        _$("#"+paramDivId).find("input,select,textarea").each(function(i){
            var jqObj = _$(this);
            var name = jqObj.attr("name");
            if(!name || !params[name]){
                return true;
            }
            var paramsVal = params[_$.trim(name)];
            var type = jqObj.attr("type");
            if(type == "text"  || type == "hidden" || jqObj.is("textarea") || jqObj.is("select") ){
                jqObj.val(paramsVal);
            }
            if(type == "radio"){
                var radioval = jqObj.val();
                jqObj.attr("checked",(radioval+"") == (paramsVal));
            }
            if(type == "checkbox"){

                var checkval = jqObj.val();
                if(_$.isArray(paramsVal) &&  _$.inArray(checkval,paramsVal) ){
                    jqObj.attr("checked",true);
                }else if(_$.type(paramsVal) == "string"){
                    jqObj.attr("checked",(checkval == paramsVal));
                }

            }
        });
    },
    delParentTag:function(thisObj,tag){
        var jqObj = yihg_utils_fun.getParentTag(thisObj,tag);
        if(jqObj){
            jqObj.remove();
        }
    },
    getParentTag:function(thisObj,tag){
        var _$ = jQuery || $;
        var obj = _$(thisObj).parent();
        while(!obj.is(tag) && !obj.is('body')){
            obj =  obj.parent()
        }
        if(obj.is(tag)){
            return obj;
        }
        return null;
    },
    getCtx:function(url){
        var newUrl = url.replace(/^\//,"");
        return yihg_erp_web_config["ctxPath"]+"/"+newUrl;
    },
    crosUnionArray : function(orginArray,dstArray,delKey,addKey){
        var delname = delKey || "delIds";
        var addname = addKey || "addIds";
        var retMap = {};
        retMap[delname] = [];
        retMap[addname] = [];
        var orginMap = {};
        var dstMap = {};

        $.each(orginArray || [],function(i,o){
            orginMap[o+""] = 1;
        });

        $.each(dstArray || [],function(i,o){
            dstMap[o+""] = 1;
        });
        //获取删除的id集合
        $.each(orginArray || [],function(i,o){
            if(!dstMap[o+""]){
                (retMap[delname]).push(o+"");
            }
        });

        //获取新增的id集合
        $.each(dstArray,function(i,o){
            if(!orginMap[o+""]){
                (retMap[addname]).push(o+"");
            }
        });

        return retMap;
    },
    parseArrayToBooleanMap:function(array,boolean){
        var a = array || [];
        var bf = boolean || 1;
        var retMap = {};
        $.each(a,function(index,o){
            retMap[o+""] = bf;
        });
        return retMap;
    },
    getDateStr:function(dateStr,addOrMinusDay){
        if(dateStr && addOrMinusDay){
            var curDate = new Date(dateStr);
            var newDate = new Date(curDate.getTime()+addOrMinusDay*86400000);
            var fullYear = newDate.getFullYear();
            var month = newDate.getMonth()+1;
            var day = newDate.getDate();
            var monthStr = "";
            var dayStr  = "";
            if(month < 10){
                monthStr = "0"+month;
            }else{
                monthStr = ""+month;
            }
            if(day < 10){
                dayStr = "0"+day;
            }else{
                dayStr = ""+day;
            }
            return fullYear+"-"+monthStr+"-"+dayStr;
        }
        return dateStr;
    },
    removeTdTitle:function(tableid,indexs){

        var tableObj = $("#"+tableid);
        tableObj.find("tr").each(function(i,o){
            var trObj = $(o);
            $(indexs).each(function(i,tindex){
                trObj.find("td:eq("+tindex+")").attr("title","");
            });
        });


    }
}


function YihgMergeTable(options){
    this.config = {
        sortIndex:0,
        tableId:"",
        cascadeMergeIndex:[],
        orderBy:1
    };
    this.cache ={
        trArray:[]
    };
    this.isSort = 0;
    this.init(options);
}

YihgMergeTable.prototype = {
    constructor:YihgMergeTable,
    init:function(options){
        this.config = $.extend(this.config,options || {});
    },
    sortTable:function(){
        var self = this,_config = self.config,_cache = self.cache;
        if(_config.tableId){
            var tbody = $("#"+_config.tableId).find("tbody")[0];
            for(var i = 0,iLength = tbody.rows.length;i<iLength;i++){
                _cache.trArray[i] = tbody.rows[i];
            }
        }
        _cache.trArray.sort(function(tr1,tr2){
            return (tr1.cells[_config.sortIndex].innerHTML+"").localeCompare((tr2.cells[_config.sortIndex].innerHTML+"")) * _config.orderBy;
        });

        //把排序后的tr 重新插入tbody
        for(var j =0; j < _cache.trArray.length; j++ ) {
            tbody.appendChild(_cache.trArray[j]);
        }
        this.isSort = 1;
    },
    show:function(){
        console.info(this.config);
        console.info(this.isSort);
    },
    mergeTable:function(){
        //根据排序列进行合并
        var self = this,_config = self.config,_cache = self.cache;
        if(!self.isSort){
            self.sortTable();
        }
        var tdIndex = _config.sortIndex;

        var cascadeTdIndexArray = _config.cascadeMergeIndex;
        //合并单元格
        var curTrIndex = 0,rowspanLenth = 1,curTdText = "";
        if(!_cache.trArray.length){
            return;
        }
        curTdText = _cache.trArray[curTrIndex].cells[tdIndex].innerHTML+"";
        var jLength = _cache.trArray.length;
        var mMap = {};
        for(var j = 0; j < jLength;j++){
            if(j+1 < jLength){
                if(!curTdText.localeCompare( _cache.trArray[(j+1)].cells[tdIndex].innerHTML+"" )){
                    $(_cache.trArray[j+1].cells[tdIndex]).css("display","none");
                    $.each(cascadeTdIndexArray,function(i,o){
                        (o+"") != (tdIndex+"") && $(_cache.trArray[j+1].cells[o]).css("display","none");

                    });
                    rowspanLenth++;
                }else{
                    $(_cache.trArray[curTrIndex].cells[tdIndex]).attr("rowspan",rowspanLenth);
                    $.each(cascadeTdIndexArray,function(i,o){
                        (o+"") != (tdIndex+"") && $(_cache.trArray[curTrIndex].cells[o]).attr("rowspan",rowspanLenth);
                    });

                    rowspanLenth = 1;
                    curTrIndex = j+1;
                    curTdText = _cache.trArray[curTrIndex].cells[tdIndex].innerHTML+"";
                }
            }else{
                $(_cache.trArray[curTrIndex].cells[tdIndex]).attr("rowspan",rowspanLenth);
                $.each(cascadeTdIndexArray,function(i,o){
                    (o+"") != (tdIndex+"") && $(_cache.trArray[curTrIndex].cells[o]).attr("rowspan",rowspanLenth);
                });
            }
        }

    }
}



/**
 * 产品品牌弹出窗
 */
function YMCommonComponent(options){
    this.config  = {
        okFunc:0,
        closeFunc:0,
        title:"选择产品.品牌",
        multi:1,
        area: ['800px', '520px']
    };
    this.init(options);
}
YMCommonComponent.prototype = {
    constructor:YMCommonComponent,
    init:function(options){
        this.config = $.extend(this.config,options || {});
    },
    openBrand:function(){
        var _this = this;
        var iframeWin ;
        var formUtils = yihg_utils_fun;
        var url = formUtils.getCtx("/component/queryProductBrandList.htm?multi="+_this.config.multi);
        layer.open({
            type:2,
            title:_this.config.title || "请选择",
            area:_this.config.area,
            content:url,
            btn: ['确定', '取消'],
            success:function(layero,index){
                //保存layer的index
                _this.config["layerIndex"] = index;
                iframeWin = window[layero.find('iframe')[0]['name']];
            },
            yes:function( index,layero){
                //获取已经选择的品牌数据对象
                var listDatas = iframeWin.getListDatas() || [];
                //回调函数
                if((typeof _this.config.okFunc) == "function"){
                    _this.config.okFunc.call({},index,listDatas || []);
                }
            },
            cancel:function(layero,index){
                layer.close(index);
            }
        });
        return _this;
    }
}


function B2BTabOptComponent(options){
    this.config = {
        title:"操作",
        menus:[{
            name:"编辑",
            click:function(rowParams){
            },
            text:function(rowParams){
                return "";
            },
            show:function(rowParams){
                return true;
            }
        }],
        keys:[]
    };
    this.init(options);
}

B2BTabOptComponent.prototype = {
    constructor:B2BTabOptComponent,
    init:function(options){
        this.config = $.extend(this.config,options || {});
        var _thisO = this;
        $("#tableDiv").on("click",".b2b-tab-opt-a",function(e){
            var jqO = $(this);
            var index = jqO.attr("aindex");
            //获取参数
            var p = {};
            $.each(_thisO.config.keys,function(i,o){
                p[o] =  jqO.attr(o);
            });
            //根据index  获取绑定的click事件
            var callBackFun = (_thisO.config.menus[index]).click;
            callBackFun && ( (typeof callBackFun )== "function") && callBackFun.call({},p);
        });
    },
    showText:function(rowParams){
        //解析文本
        var _this = this;
        //获取属性
        if(!_this.config.keys.length && rowParams && (rowParams instanceof Object)){
            $.each(rowParams,function(key,value){
                _this.config.keys.push(key);
            });
        }
        var topDivJq = $("<div></div>")
        //显示标题
        var  titleJq = $("<div class='tab-operate'><a href='javascript:void(0);' class='btn-show'>"+_this.config.title+"<span class='caret'></span></a></div>")
        if(_this.config.menus && _this.config.menus.length){
            var btnHiddeJq = $("<div class='btn-hide'></div>")
            $(_this.config.menus).each(function(i,o){

                if(o.show && (typeof o.show) == "function" ){
                    var isShow = o.show(rowParams);
                    if(!isShow){
                        //不显示该菜单项
                        return true;
                    }
                }

                var menuName = "";
                if(o.text && (typeof o.text) == "function"){
                    menuName  = o.text(rowParams);
                }else{
                    menuName = o.name || "##";
                }
                var menuJq = $("<a href='javascript:void(0)' aindex='"+i+"' class='def b2b-tab-opt-a'>"+menuName+"</a>");
                //用属性存参数
                _this.config.keys.length && $.each(rowParams,function(key,value){
                    menuJq.attr(key,value);
                });
                btnHiddeJq.append(menuJq);
            });
            titleJq.append(btnHiddeJq);
        }
        topDivJq.append(titleJq);
        return topDivJq.html();
    }
}


function B2BCheckComponent(options){
    this.config = {
        okFunc:0,
        closeFunc:0,
        title:"请选择",
        url:"",
        area: ['900px', '650px'],
        maxmin:true,
        params:{},
        shadeClose : true,
        shade : 0.5,
        btn: ['确定', '取消'],
        pfun:"_getAllDatas",
        pwin:""
    };
    this.init(options);
}
B2BCheckComponent.prototype = {
    constructor:B2BCheckComponent,
    init:function(options){
        this.config = $.extend(this.config,options || {});
    },
    _getNewUrl:function(){
        var formUtils = yihg_utils_fun;
        var _this = this;
        var _url =_this.config.url;
        if(_this.config.params){
            var ps = "";
            for(var _key in _this.config.params ){
                ps += ("&"+_key)+"="+_this.config.params[_key];
            }
            if(_url.indexOf("?") <= 0){
                _url += "?_t=123";
            }
            _url += ps;
        }
        return encodeURI(_url);
    },
    setParams:function(params){
        var _this = this;
        _this.config.params = params;
        return _this;
    },
    open:function(e){
        var  _this = this;
        var iframeWin ;
        var newUrl = _this._getNewUrl();
        layer.open({
            type:2,
            title:_this.config.title || "请选择",
            area:_this.config.area,
            content:newUrl,
            btn: _this.config.btn,
            shadeClose : _this.config.shadeClose,
            shade : _this.config.shade,
            success:function(layero,index){
                //保存layer的index
                _this.config["layerIndex"] = index;
                iframeWin = window[layero.find('iframe')[0]['name']];
            },
            yes:function( index,layero){
                //获取已经选择的品牌数据对象
                //var listDatas = iframeWin._getAllDatas() || {};
                var listDatas = ((iframeWin[_this.config.pfun])()) || {};
                //回调函数
                if((typeof _this.config.okFunc) == "function"){
                    _this.config.okFunc.call({},index,listDatas,e);
                }
            },cancel: function(index){
                layer.close(index);
            },
            success:function(layero,index){
                iframeWin = window[layero.find('iframe')[0]['name']];
                _this.config["pwin"] = iframeWin;
                _this.config["idx"] = index;
            }
        });
        return _this;
    },
    getDatas:function(){
        var _this = this;
        //在弹窗未关闭之前获取前一个页面的返回的数据
        var pwin = _this.config["pwin"];
        return ((pwin[_this.config.pfun])()) || {};
    },
    close:function(){
        var _this = this;
        layer && _this.config["idx"] && layer.close(_this.config["idx"]);
        return  _this;
    }
}






