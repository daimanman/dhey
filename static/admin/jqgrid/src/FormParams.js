layui.define(function(exports) {
	"use strict";
	var $ = layui.jquery;
	var obj = {
		getParams : function(paramDivId) {
			var params = {};
            $("#"+paramDivId).find(" input , select ").each(function(i){
                var jqObj = $(this);
                var type = jqObj.attr("type");
                var val = jqObj.val();
                var name = jqObj.attr("name");

                if(type == "text" || type=="hidden" || jqObj.is("select")){
                    params[name] = val ;
                }

                if(type == "radio")
                {
                    jqObj.attr("checked") && (params[name] = val);
                }

                if(type =="checkbox" ){
                    var checkAttr = jqObj.attr("checked");

                    if(checkAttr && !params[name]){
                        params[name] = [];
                    }

                    checkAttr && (params[name]).push(val);
                }
            });
           return params;
		}
	};
	// 输出test接口
	exports('FormParams', obj);
});