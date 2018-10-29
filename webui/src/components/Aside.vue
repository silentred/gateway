<template>
    <aside class="main-sidebar">
        <!-- sidebar: style can be found in sidebar.less -->
        <section class="sidebar">
            <!-- Sidebar Menu -->
            <ul class="sidebar-menu">
                <li class="header">后台目录</li>
                <li class="treeview active">
                    <a href="#">
                        <i class="fa fa-table"></i>
                        <span>路由</span>
                    </a>
                    <ul class="treeview-menu">
                        <li><router-link :to="{ name: 'RouteList'}"><i class="fa fa-circle-o"></i>路由列表</router-link></li>
                        <li><router-link :to="{ name: 'RouteForm'}"><i class="fa fa-circle-o"></i>新建路由</router-link></li>
                    </ul>
                </li>
                <li class="treeview active">
                    <a href="#">
                        <i class="fa fa-table"></i>
                        <span>后端服务</span>
                    </a>
                    <ul class="treeview-menu">
                        <li><router-link :to="{ name: 'ServiceList'}"><i class="fa fa-circle-o"></i>服务列表</router-link></li>
                    </ul>
                </li>
                <!-- Optionally, you can add icons to the links -->
                <li class="treeview active">
                    <a href="#"><i class="fa fa-link"></i> <span>后台管理</span></a>
                    <ul class="treeview-menu">
                        
                    </ul>
                </li>
            </ul>
            <!-- /.sidebar-menu -->
        </section>
        <!-- /.sidebar -->
    </aside>
</template>

<script>
export default {
    name: 'Aside',
    data () {
        return {}
    }
}
</script>