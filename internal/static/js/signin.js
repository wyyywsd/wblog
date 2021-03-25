// $.validator.setDefaults({
//     submitHandler: function() {
//         alert("提交事件!");
//
//     }
// });
$().ready(function () {
    // 在键盘按下并释放及提交后验证提交表单
    $("#signinForm").validate({
        rules: {
            username: {
                required: true,
                minlength: 2
            },
            password: {
                required: true,
                minlength: 5
            }

        },
        messages: {
            username: {
                required: "请输入用户名",
                minlength: "用户名必需由两个字母组成"
            },
            password: {
                required: "请输入密码",
                minlength: "密码长度不能小于 5 个字母"
            }
        }
    });
});