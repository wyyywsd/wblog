$(document).ready(function () {
    if (unbind_status == "pending") {
        $("#pending").attr("selected", "selected");
    } else if (unbind_status == "processing") {
        $("#processing").attr("selected", "selected");
    } else if (unbind_status == "success") {
        $("#success").attr("selected", "selected");
    }

    $('#exampleModal').on('show.bs.modal', function (event) {
        var button = $(event.relatedTarget) // 触发事件的按钮
    })
    $('#exampleModal2').on('show.bs.modal', function (event) {
        var button = $(event.relatedTarget) // 触发事件的按钮
    })

    $("#testFile").change(function () {
        run(this, function (data) {
            alert(data)
            $('#testImg').attr('src', data);
            $('#testArea').val(data);
        });
    });
// $("#btnTest").click(function () {
//     $.ajax({
//         url: "/usercenter/testbaseaction",
//         type: "post",
//         dataType: "json",
//         data: {
//             "content": $("#testArea").val(),
//         },
//         async: false,
//         success: function (result) {
//             if (result.Code == 200) {
//                 alert(result.Data);
//             } else {
//             }
//         }
//     });
// });

    function run(input_file, get_data) {
        /*input_file：文件按钮对象*/
        /*get_data: 转换成功后执行的方法*/
        if (typeof (FileReader) === 'undefined') {
            alert("抱歉，你的浏览器不支持 FileReader，不能将图片转换为Base64，请使用现代浏览器操作！");
        } else {
            try {
                /*图片转Base64 核心代码*/
                var file = input_file.files[0];
                //这里我们判断下类型如果不是图片就返回 去掉就可以上传任意文件
                if (!/image\/\w+/.test(file.type)) {
                    alert("请确保文件为图像类型");
                    return false;
                }
                var reader = new FileReader();
                reader.onload = function () {
                    get_data(this.result);
                }
                reader.readAsDataURL(file);
            } catch (e) {
                alert('图片转Base64出错啦！' + e.toString())
            }
        }
    }


});


