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
            jqGridTabObj.setGridWidth(width);
            jqGridTabObj.setGridHeight(height+70);
        },
        loadGrid: function () {

            var  URL= _webConfig.api+_webConfig.uinfo_list;
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
                    {name: 'id',index: 'id',width: 20, align: "center", sortable: false,label:"ID"},
                    {name: 'icon',index: 'id',width: 30, align: "center", sortable: false,label:"ICON",formatter:formatFun.getQIcon},
                    {name: 'sex',index: 'id',width: 30, align: "center", sortable: false,label:"性别",formatter:formatFun.getSexIcon},
                    {name: 'uid',index: 'uid',width: 100, align: "center", sortable: false,label:"UID"},
                    {name: 'nickname',index: 'nickname',width: 80, align: "center", sortable: false,label:"昵称"},
                    {name: 'age',index: 'age',width: 30, align: "center", sortable: false,label:"年龄"},
                    {name: 'career',index: 'career',width: 80, align: "center", sortable: false,label:"职业"},
                    {name: 'country',index: 'country',width: 50, align: "center", sortable: false,label:"现居国"},
                    {name: 'province',index: 'province',width: 50, align: "center", sortable: false,label:"现居省"},
                    {name: 'city',index: 'city',width: 50, align: "center", sortable: false,label:"现居城市"},
                    {name: 'hco',index: 'hco',width: 50, align: "center", sortable: false,label:"家国"},
                    {name: 'hp',index: 'uid',width:50, align: "center", sortable: false,label:"家省"},
                    {name: 'hc',index: 'uid',width: 50, align: "center", sortable: false,label:"家城"},
                    {name: 'marriage',index: 'marriage',width:30, align: "center", sortable: false,label:"婚姻"},
                    {name: 'company',index: 'company',width: 80, align: "center", sortable: false,label:"公司"},
                    {name: 'birthyear',index: 'birthyear',width: 50, align: "center", sortable: false,label:"出生年"},
                    {name: 'birthday',index: 'birthday',width: 50, align: "center", sortable: false,label:"出生月日"},
                    {name: 'bloodtype',index: 'bloodtype',width: 30, align: "center", sortable: false,label:"血型"}
                ],
                regional : 'cn',
                pager: "#pagerDiv",
                viewrecords: true,
                caption: "",
                multiselect:true,
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
