function login() {
    var email = document.getElementById("email").value;
    var passwd = document.getElementById("passwd").value;
    if(document.getElementById("login_day").checked){
        var login_day = "0"
    }else {
        var login_day = "1"
    }
    $.ajax(
        {
            type: 'post',
            url: '/login',
            data: {"email":  email, "passwd": passwd, "login_day" : login_day},
            dataType: 'json',
            success: function (data) {
                if (data.code != 1) {
                    layer.msg(data.msg)
                }else {
                    window.location.href = "/";
                }
            }
        }
    );
}
