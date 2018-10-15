function _refresh(avg) {
    layer.msg('数据加载中...', {icon: 16, shade: 0.01, time: 0});
    $.ajax({
        type: 'post', url: '/ldap', data: {"action": avg}, dataType: 'json', success: function (data) {
            if (data.code == 0) {
                var trStr = '';
                if (avg == "LdapUserInfo") {
                    if (data.data.length == 0) {
                        layer.msg("未查询到数据", {icon: 5});
                    } else {
                        for (i = 0; i < data.data.length; i++) {
                            trStr += '<tr> ';
                            trStr += '<td><p style="font-weight:bold">' + data.data[i].uid + '</p></td> <td>' + data.data[i].displayName + '</td> <td>' + data.data[i].mail + '</td> <td>' + data.data[i].mobile + '</td><td>' + data.data[i].group + '</td> ';
                            trStr += '<td><a class="layui-btn layui-btn-radius layui-btn-xs layui-btn-primary"> 详细 </a>' + '<a class="layui-btn layui-btn-radius layui-btn-xs"> 编辑 </a>' + '<a class="layui-btn layui-btn-radius layui-btn-xs layui-btn-danger"> 删除 </a>';
                        }
                        $("#LdapUserInfo").html(trStr);
                        layer.msg("数据加载完成")
                    }
                } else if (avg == "LdapGroupInfo") {
                    if (data.data.length == 0) {
                        layer.msg("未查询到数据", {icon: 5});
                    } else {
                        for (i = 0; i < data.data.length; i++) {
                            trStr += '<tr> ';
                            trStr += '<td><p style="font-weight:bold">' + data.data[i].cn + '</p></td> <td>' + data.data[i].gidNumber + '</td> <td>' + data.data[i].description + '</td>';
                            trStr += '<td><a class="layui-btn layui-btn-radius layui-btn-xs layui-btn-primary"> 详细 </a>' + '<a class="layui-btn layui-btn-radius layui-btn-xs"> 编辑 </a>' + '<a onclick=DelGroup("' + data.data[i].cn + '") class="layui-btn layui-btn-radius layui-btn-xs layui-btn-danger"> 删除 </a>';
                        }
                        $("#LdapGroupInfo").html(trStr);
                        layer.msg("数据加载完成")
                    }
                }

            } else if (data.code == -1) {
                layer.msg(data.msg, {icon: 5});
            } else {
                layer.alert(data.msg, {
                    skin: 'layui-layer-molv' //样式类名
                    , closeBtn: 0
                }, function () {
                    window.location.href = "/login";
                });
            }
        }
    });
}


function DelGroup(cn) {
    layer.alert("确定删除该用户组? ", {
        skin: 'layui-layer-molv' //样式类名
        , closeBtn: 0
    }, function () {
        $.ajax(
            {
                type: 'post',
                url: '/ldap',
                data: {"cn": cn, "action": 'DleGroup'},
                dataType: 'json',
                success: function (data) {
                    if (data.code != 1) {
                        layer.msg(data.msg)
                    } else {
                        parent.layer.closeAll();
                        _refresh("LdapGroupInfo")
                    }
                }
            }
        )
    });

}

function _Grouprefresh() {
    $("#LdapUserDiv").hide();
    $("#LdapGroupDiv").show();
    _refresh("LdapGroupInfo")

}

function _Userrefresh() {
    $("#LdapUserDiv").show();
    $("#LdapGroupDiv").hide();
    _refresh("LdapUserInfo")

}