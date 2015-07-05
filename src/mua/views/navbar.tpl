{{define "navbar"}}

<a class="navbar-brand" href="/">MUA</a>
<div>
	<ul class="nav navbar-nav">
	    <li {{if .IsHome}}class="active"{{end}}><a href="/">首页</a></li>
	    <li {{if .IsMy}}class="active"{{end}}><a href="/my">上传表情</a></li>
	</ul>
</div>

<div class="pull-right">
	<ul class="nav navbar-nav">
		{{if .IsLogin}}
		<li><a href="/login?exit=true">退出</a></li>
		{{else}}
		<li><a href="/login">管理员登录</a></li>
		{{end}}
	</ul>
</div>
{{end}}