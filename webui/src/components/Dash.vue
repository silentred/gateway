<template>
    <div class="wrapper">
        <top-header></top-header>
        <left-aside></left-aside>
        <div class="content-wrapper">
            <router-view></router-view>
        </div>
        <footer class="main-footer">
            <strong>Copyright &copy; 2016 <a href="http://fameapp.us">fameapp.us</a> </strong> All rights reserved.
        </footer>
    </div>
</template>

<script>
import Header from '@/components/Header'
import Aside from '@/components/Aside'

export default {
    name: 'Dash',
    data () {
        return {}
    },
    components: {
        'top-header': Header,
        'left-aside': Aside
    }
}
</script>