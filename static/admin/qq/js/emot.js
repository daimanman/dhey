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
            jqGridTabObj.setGridWidth(width-20);
            jqGridTabObj.setGridHeight(height+70);
        },
        loadGrid: function () {

            var  URL= _webConfig.api+_webConfig.emot_list;
            jqGridTabObj.jqGrid({
                url: URL,
                datatype: "json",
                mtype : "post",
                height: 500,
                autowidth: false,
                shrinkToFit: true,
                rownumbers:true,
                rowNum: 15,
                rowList: [10,15, 30, 50, 100, 500, 1000],
                colModel: [
                    {name: 'id',index: 'id',sortable: true,width: 30, align: "center",label:"ID"},
                    {name: 'content',index: 'content',width: 220, align: "center", sortable: false,label:"content"},
                    {name: 'createTime',index: 'createTime',width: 60, align: "center", sortable: false,label:"时间"},
                    {name: 'uid',index: 'uid',width: 40, align: "center", sortable: false,label:"UID"},
                    {name: 'name',index: 'name',width: 80, align: "center", sortable: false,label:"name"},
                    {name: 'pictotal',index: 'pictotal',width: 20, align: "center", sortable: true,label:"图片"},
                    {name: 'lbs_name',index: 'lbs_name',width: 80, align: "center", sortable: false,label:"地点"},
                    {name: 'source_name',index: 'source_name',width: 50, align: "center", sortable: false,label:"手机"}

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
                },
                loadComplete:function(data){


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
