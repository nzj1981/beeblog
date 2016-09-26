<div class="bs-example bs-navbar-top-example">
    <nav role="navigation" class="navbar navbar-default navbar-fixed-top">
      <!-- We use the fluid option here to avoid overriding the fixed width of a normal container within the narrow content columns. -->
      <div class="container-fluid">
        <div class="navbar-header">
          <button data-target="#bs-example-navbar-collapse-6" data-toggle="collapse" class="navbar-toggle collapsed" type="button">
            <span class="sr-only">导航</span>
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
          </button>
          <a href="/" class="navbar-brand">我的博客</a>
        </div>

        <div id="bs-example-navbar-collapse-6" class="collapse navbar-collapse">
          <ul class="nav navbar-nav">
            <li {{if .IsHome}}class="active"{{end}}><a href="/">首页</a></li>
            <li {{if .IsCategory}}class="active"{{end}}><a href="/category">分类</a></li>
            <li {{if .IsTopic}}class="active"{{end}}><a href="/topic">文章</a></li>
          </ul>
  <div class="pull-right">
    <ul class="nav navbar-nav">
    {{if .IsLogin}}
    <li><a href="/login?exit=true">退出</a></li>
    {{else}}
    <li><a href="/login">管理员登录</a></li>
    </ul>
    {{end}}
  </div>
  </div><!-- /.navbar-collapse -->
  </div>
  </nav>
</div>