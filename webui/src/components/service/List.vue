<template>
    <div class="hello">
        <section class="content-header">
            <h1>后端服务</h1>
            <ol class="breadcrumb">
                <li>
                    <a href="#">
                        <i class="fa fa-dashboard"></i> 管理后台</a>
                </li>
                <li class="active">后端服务</li>
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
                                    <span class="clickable text-danger" :target-id="val.id" :target-host="val.host" :svc-name="item.item.service.name"  :data-index="item.index" @click="delTarget">
                                        <i class="fa fa-fw fa-remove"></i>
                                        删除
                                    </span>
                                </li>
                            </ul>
                        </template>
                        <template slot="guards" scope="item">
                            {{item.item.service.guards.name}}
                        </template>
                        <template slot="reactors" scope="item">
                            {{item.item.service.reactors.name}}
                        </template>
                        <template slot="actions" scope="item">
                            <b-button variant="primary" size="sm" :data-index="item.index" @click="clickAddTarget">添加Target</b-button>
                        </template>
                    </b-table>
                </div>
                <!-- /.box-body -->
                <div class="box-footer clearfix">
                    <b-pagination size="md" :total-rows="this.items.length" :per-page="perPage" v-model="currentPage" />
                </div>

                <div id="modal">
                    <!-- Modal Component -->
                    <b-modal id="addTargetModal" hideHeaderClose title="添加Target" @ok="submitModal" @shown="modalShowing">
                        <div role="form">
					        <div id="formWrapper">
                                <b-form-fieldset description="target host" label="Target Host">
                                    <b-form-input v-model="target.host"></b-form-input>
                                </b-form-fieldset>
                                <b-form-fieldset description="target port" label="Target Port">
                                    <b-form-input v-model="target.port"></b-form-input>
                                </b-form-fieldset>
                                <b-form-fieldset description="target weight" label="Target Weight">
                                    <b-form-input v-model="target.weight"></b-form-input>
                                </b-form-fieldset>
                                <b-form-fieldset description="health check address" label="Health Check Address">
                                    <b-form-input v-model="target.hc_addr"></b-form-input>
                                </b-form-fieldset>
                                <b-form-fieldset description="health check interval" label="Health Check Interval">
                                    <b-form-input v-model="target.hc_interval"></b-form-input>
                                </b-form-fieldset>
					        </div>
                        </div>
                    </b-modal>
                </div>
            </div>
        </section>
    </div>
</template>

<script>
import 'bootstrap-vue/dist/bootstrap-vue.css'
import { bTable, bPagination, bFormFieldset, bFormInput, bButton } from 'bootstrap-vue/lib/components'
import bModal from '@/components/bootstrap/Modal'
import axios from 'axios'
import config from '@/config'

export default {
    name: 'ServiceList',
    components: {
        bTable, bPagination, bFormFieldset, bFormInput, bButton, bModal
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
                },
                actions: {
                    label: '操作',
                }
            },
            currentPage: 1,
            perPage: 15,
            filter: null,
            target: {
                host: '',
                port: 80,
                weight: 10,
                hc_addr: '',
                hc_interval: '10s',
                index: 0
            }
        }
    },
    methods: {
        delTarget(event){
            var that = this;
            var index = event.target.getAttribute("data-index");
            var svcID = event.target.getAttribute("target-id");
            var targetHost = event.target.getAttribute("target-host");
            var svcName = event.target.getAttribute("svc-name");
            
            console.log(svcID, index);
            var path = '/table/target/' + svcID;
            var domain = window.location.hostname;
        	var url = config.scheme + '://' + domain + ':' + config.port + path; 

            axios.delete(url, {
                headers: {'Content-Type': 'application/json'},
                params: {
                    service_name:svcName,
                    target_host:targetHost
                }
            }).then(function(response){
                if (response.data.code == 0) {
                    // remove target from UI
                    //that.items.splice(index, 1)
                    //window.location.reload();
                    alert("ok");
                }else {
                    alert(response.data.msg);
                }
                console.log(response.data)
            }).catch(e => {
                console.log(e)
            })

            return false
        },
        modalShowing(){
            console.log(this.target)
        },
        submitModal(){
            var item = this.items[this.target.index]
            console.log(item)

            var obj = {
                route_host: item.route.host,
                route_prefix: item.route.prefix,
                service_strip: item.service.strip,
                service_name: item.service.name,
                target_host: this.target.host,
                target_port: parseInt(this.target.port),
                target_weight: parseInt(this.target.weight),
                hc_addr: this.target.hc_addr,
                hc_interval: this.target.hc_interval
            };
            console.log("submitting:", obj)

            var path = '/table/service';
            var domain = window.location.hostname;
        	var url = config.scheme + '://' + domain + ':' + config.port + path;

			axios.put(url, obj, {
				headers: {'Content-Type': 'application/json'}
			}).then(function(response){
				console.log("then", response)
			}).catch(e => {
				console.log("catching", e)
			})
        },
        clickAddTarget(event){
            var index = event.target.getAttribute("data-index");
            this.target.index = parseInt(index)
            this.$root.$emit('show::modal','addTargetModal')
        }
    },
    watch: {
        'target.host': function(){
            this.target.hc_addr = this.target.host + ':' + this.target.port
        },
        'target.port': function(){
            this.target.hc_addr = this.target.host + ':' + this.target.port
        }
    }
}
</script>

<style scoped>
.clickable {
    cursor: pointer;    
}
</style>
