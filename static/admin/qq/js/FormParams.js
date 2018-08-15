/**

 * daixiaoman  获取表单参数
 */
yihg_utils_fun = {
	getParams : function(paramDivId) {
		var params = {};
		var _$ = jQuery || $;
		_$("#" + paramDivId).find(" input , select ,textarea").each(
				function(i) {
					var jqObj = _$(this);
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
    }
    ,
	getQIcon:function(uid){
    	return "<img src='http://qlogo1.store.qq.com/qzone/"+uid+"/"+uid+"/30' />"
	},
	parseUrlParams:function(href){
		var p = {};
		if(href && href.length > 0){
			var firstSplits = href.split("?");
			if(firstSplits.length){
				var secondSplits = firstSplits[1].split("&");
				for(var i = 0;i<secondSplits.length;i++){
					var thirdSplits = secondSplits[i].split("=");
					if(thirdSplits.length > 1){
						p[thirdSplits[0]] = thirdSplits[1];
					}else{
						p[thirdSplits[0]] = "";
					}
				}
			}
		}
		return p;
	}
}












