layui.use(['form'],function(){
    var $ = layui.jquery,
        layer = layui.layer,
        formUtils = yihg_utils_fun,
        jqGridTabObj = $("#tableDiv");
    var form = layui.form();
    var formatFun = {
        getQIcon:function(cellVal,option,rowObj){
            return formUtils.getQIcon(rowObj["uid"]);
        },
        getSexIcon:function(cellVal,option,rowObj){
            var icon = {"1":"<icon class='fa fa-male fa-2x ' style='color:blue;'/>","2":"<icon class='fa fa-female  fa-2x' style='color:indianred;' />"};
            if(rowObj["sex"]){
                return icon[rowObj["sex"]+""] || "-";
            }
            return "<icon class='fa fa-meh-o  fa-2x'/>";
    }
    }
    var opGrid = {
        reSize: function () {
            var width = $('.jqGrid_wrapper').width();
            var height = $(window).height();//parent.get_MainContainerHeight();
            var searchBox = 120, jqGrid_head = 55, jqGrid_pager = 30, jqGrid_footer = 45;
            height = height - searchBox - jqGrid_head - jqGrid_pager - jqGrid_footer;
            jqGridTabObj.setGridWidth(width-30);
            jqGridTabObj.setGridHeight(height+70);
        },
        loadGrid: function () {

            var  URL= _webConfig.api+_webConfig.phot_list;
            jqGridTabObj.jqGrid({
                url: URL,
                datatype: "json",
                mtype : "post",
                height: 500,
                autowidth: false,
                shrinkToFit: true,
                rownumbers:true,
                rowNum: 10,
                rowList: [10,15, 30, 50, 100, 500, 1000],
                colModel: [
                    {name: 'id',index: 'id',width: 50, align: "center", sortable: false,label:"ID"},
                    {name: 'pre',index: 'pre',width: 60, align: "center", sortable: false,label:"pre",formatter:function(pre,o,rowObj){
                        return "<img width='100' height='100' src='"+pre+"'/>";
                    }},
                    {name: 'name',index: 'name',width: 60, align: "center", sortable: false,label:"name"},
                    {name: 'totalnum',index: 'totalnum',width: 30, align: "center", sortable: false,label:"totalnum"},
                    {name: 'allow_access',index: 'allow_access',width: 30, align: "center", sortable: false,label:"ICON"},
                    {name: 'comment',index: 'comment',width: 30, align: "center", sortable: false,label:"comment"},
                    {name: 'uid',index: 'uid',width: 60, align: "center", sortable: false,label:"UID"},
                    {name: 'desc',index: 'desc',width: 180, align: "center", sortable: false,label:"desc"},

                    {name: 'createtime',index: 'createtime',width: 60, align: "center", sortable: false,label:"time",formatter:function(time,opt,rowObj){
                        var date = new Date(time*1000);
                        return date.getFullYear()+"-"+(date.getMonth()+1)+"-"+date.getDate();
                    }}

                ],
                regional : 'cn',
                pager: "#pagerDiv",
                viewrecords: true,
                caption: "",
                multiselect:false,
                postData:formUtils.getParams('admin-main-params'),
                jsonReader:{
                    id:"id",
                    root:"datas",
                    page:"page",
                    total:"totalPage",
                    records:"total"
                }
            });
        }
    };
   opGrid.loadGrid();
    opGrid.reSize();
    $(window).bind('resize', function () {
        opGrid.reSize();
    });

    $('#searchBtn').on('click', function () {
        jqGridTabObj.jqGrid('setGridParam',{
            postData:formUtils.getParams('admin-main-params'),
            page:1
        }).trigger("reloadGrid");
    });

});
