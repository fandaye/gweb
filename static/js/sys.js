function _SysUserrefresh() {$("#SysUserDiv").show();$("#SysMenuDiv").hide();$("#SysRoleDiv").hide();_Sysrefresh("SysUserInfo")}function _SysMenurefresh() {$("#SysUserDiv").hide();$("#SysMenuDiv").show();$("#SysRoleDiv").hide();_Sysrefresh("SysMenuInfo")}function _SysRolerefresh() {$("#SysUserDiv").hide();$("#SysMenuDiv").hide();$("#SysRoleDiv").show();_Sysrefresh("SysRoleInfo")}function _Sysrefresh(avg) {layer.msg('数据加载中...', {icon: 16, shade: 0.01, time: 0});$.ajax({type: 'post', url: '/sys', data: {"action": avg}, dataType: 'json', success: function (data) {if (data.code != '0'){puberror(data.code,data.msg);return;}var trStr = '';if (avg == "SysUserInfo") {if (data.data.length == 0) {layer.msg("未查询到数据", {icon: 5});} else {for (i = 0; i < data.data.length; i++) {trStr += '<tr> <td><p style="font-weight:bold">' + data.data[i].email + '</p></td> <td>' + data.data[i].username + '</td> <td>' + data.data[i].role_id + '</td> ' ;if (data.data[i].status == "0"){trStr += ' <td><a class="layui-btn layui-btn-radius layui-btn-xs"> 启用 </a></td><td>'}else {trStr += ' <td><a class="layui-btn layui-btn-radius layui-btn-xs layui-btn-danger"> 禁用 </a></td><td>'}trStr += data.data[i].create_time + '</td> <td>' + data.data[i].login_time + '</td>';trStr += '<td><a class="layui-btn layui-btn-radius layui-btn-xs layui-btn-primary"> 详细 </a>' + '<a class="layui-btn layui-btn-radius layui-btn-xs"> 编辑 </a>' + '<a class="layui-btn layui-btn-radius layui-btn-xs layui-btn-danger"> 删除 </a>';}$("#SysUserInfo").html(trStr);layer.closeAll()}} 
else if (avg == "SysMenuInfo") {if (data.data.length == 0) {layer.msg("未查询到数据", {icon: 5});} else {for (i = 0; i < data.data.length; i++) {trStr += '<tr> <td><p style="font-weight:bold"> <i class="lnr ' + data.data[i].menu_lcon + '"></i></p></td> <td>' + data.data[i].menu_url + '</td> <td>' + data.data[i].menu_name + '</td><td>' + data.data[i].create_time + '</td><td>' + data.data[i].instruction + '</td>';trStr += '<td><a class="layui-btn layui-btn-radius layui-btn-xs layui-btn-primary"> 详细 </a>' + '<a class="layui-btn layui-btn-radius layui-btn-xs" onclick="SysMenuEdit('+ JSON.stringify(data.data[i]).replace(/\"/g,"'") +')"> 编辑 </a>' + '<a class="layui-btn layui-btn-radius layui-btn-xs layui-btn-danger" onclick=SysMenuDel('+ data.data[i].id +')> 删除 </a>';}$("#SysMenuInfo").html(trStr);layer.closeAll()}}
else if (avg == "SysRoleInfo") {if (data.data.length == 0) {layer.msg("未查询到数据", {icon: 5});} else {for (i = 0; i < data.data.length; i++) {trStr += '<tr> <td><p style="font-weight:bold">' + data.data[i].role_name + '</p></td> <td>' + data.data[i].create_time + '</td><td>' + data.data[i].instruction + '</td>';trStr += '<td><a class="layui-btn layui-btn-radius layui-btn-xs layui-btn-primary"> 权限维护 </a>' + '<a class="layui-btn layui-btn-radius layui-btn-xs"> 编辑 </a>' + '<a class="layui-btn layui-btn-radius layui-btn-xs layui-btn-danger"> 删除 </a>';}$("#SysRoleInfo").html(trStr);layer.closeAll()}}}});}

function SysMenuDel(id) {layer.confirm('确定删除该系统菜单？', {btn: ['确定','取消'] }, function () {$.ajax({type: 'post', url: '/sys', data: {"action": 'DelMenu', "id": id}, dataType: 'json', success: function(data) {layer.msg(data.msg);if (data.code == '0'){_SysMenurefresh()}}});}, function () {layer.closeAll()});}

function SysMenuEdit(data) {
    layer.open({
        type: 1,
        title: "编辑系统菜单",
        area: ['60%', '265px'],
        skin: 'white', //没有背景色
        shadeClose: true,
        content: '<div id="EditMenuDiv"> <form class="layui-form layui-form-pane" action="" lay-filter="EditMenuDiv">  <div class="main-content"> <div class="container-fluid"> <label class="layui-form-label">地址</label> <div class="layui-input-block" > <input type="text" name="menu_url" autocomplete="off" class="layui-input" lay-verify="menu_url" value="'+ data.menu_url +'" disabled>  </div> </div> <div class="container-fluid"> <label class="layui-form-label">名称</label> <div class="layui-input-block"> <input type="text" name="menu_name" autocomplete="off" class="layui-input" lay-verify="menu_name" value="'+ data.menu_name +'"> </div> </div> <div class="container-fluid"> <label class="layui-form-label">描述</label> <div class="layui-input-block"> <input type="text" name="instruction" autocomplete="off" class="layui-input" lay-verify="instruction" value="'+ data.instruction +'"> </div> </div> <div class="container-fluid"> <div class="layui-input-block" style="float:right"> <button class="layui-btn" lay-submit="" lay-filter="EditMenuDivPost" >立即提交</button> <button type="reset" class="layui-btn layui-btn-primary">重置</button> </div> </div> </div> </form> </div>'
    });
}