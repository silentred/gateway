<template>
    <div class="hello">
        <section class="content-header">
            <h1>路由列表</h1>
            <ol class="breadcrumb">
                <li>
                    <a href="#">
                        <i class="fa fa-dashboard"></i> 管理后台</a>
                </li>
                <li class="active">路由列表</li>
            </ol>
        </section>
        <!-- Main content -->
        <section class="content">
            <div class="box">
                <div class="box-header with-border">
                    <h3 class="box-title">列表</h3>
                </div>
                <!-- /.box-header -->
                <div class="box-body table-responsive">
                    <div class="">
                        <div class="form-inline">
                            <div class="form-group">
                                <b-form-fieldset horizontal label="过滤路由" class="" :label-size="4">
                                    <b-form-input v-model="filter" placeholder="Type to Search"></b-form-input>
                                </b-form-fieldset>
                            </div>
                        </div>
                    </div>
                    <!-- Main table element -->
                    <b-table hover bordered striped :items="items" :fields="fields" :current-page="currentPage" :per-page="perPage" :filter="filter">
                        <template slot="routeID" scope="item">
                            {{item.item.route.id}}
                        </template>
                        <template slot="route" scope="item">
                            {{item.value.host}}{{item.value.prefix}}
                        </template>
                        <template slot="serviceName" scope="item">
                            {{item.item.service.name}}
                        </template>
                        <template slot="serviceStrip" scope="item">
                            {{item.item.service.strip}}
                        </template>
                        <template slot="targets" scope="item">
                            <ul>
                                <li v-for="val in item.item.service.targets.list">
                                    <span>域名:{{val.host}} 权重:{{val.weight}} </span> 
                                </li>
                            </ul>
                        </template>
                        <template slot="guards" scope="item">
                            {{item.item.service.guards.name}}
                        </template>
                        <template slot="reactors" scope="item">
                            {{item.item.service.reactors.name}}
                        </template>
                    </b-table>
                </div>
                <!-- /.box-body -->
                <div class="box-footer clearfix">
                    <b-pagination size="md" :total-rows="this.items.length" :per-page="perPage" v-model="currentPage" />
                </div>
            </div>
        </section>
    </div>
</template>

<script>
import 'bootstrap-vue/dist/bootstrap-vue.css'
import { bTable, bPagination, bFormFieldset, bFormInput, bButton } from 'bootstrap-vue/lib/components'
import axios from 'axios'
import config from '@/config'

export default {
    name: 'ServiceList',
    components: {
        bTable, bPagination, bFormFieldset, bFormInput, bButton
    },
    created: function() {
        // fetch data form server /table/routes
        var path = '/table/routes';
        var domain = window.location.hostname;
        var url = config.scheme + '://' + domain + ':' + config.port + path;
        var that = this
    
        axios.get(url).then(function(response){
            that.items = response.data
            console.log(that.items)
        }).catch(e => {
            console.log(e)
        })
    },
    data() {
        return {
            items: [],
            fields: {
                routeID: {
                    label: '路由ID',
                },
                route: {
                    label: '路由',
                    sortable: true
                },
                serviceName: {
                    label: '服务名称',
                },
                serviceStrip: {
                    label: '服务删除前缀',
                },
                targets: {
                    label: '后端服务',
                },
                guards: {
                    label: '服务过滤器',
                },
                reactors: {
                    label: '服务监控器',
                }
            },
            currentPage: 1,
            perPage: 5,
            filter: null
        }
    }
}
</script>