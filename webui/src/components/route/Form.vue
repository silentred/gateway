<template>
	<section class="content">
		<div class="box box-primary">
			<div class="box-header with-border">
				<h3 class="box-title">创建路由</h3>
			</div>
			<div class="alert" :class="[alert.class]" v-show="alert.show">
                <button type="button" class="close" @click="alert.show = !alert.show" >×</button>
                <span>{{alert.msg}}</span>
            </div>
			<!-- /.box-header -->
			<!-- form start -->
			<div role="form">
				<div class="box-body">
					<div id="formWrapper">
						<b-form-fieldset description="route host" label="Route Host" :state="routeHostState">
							<b-form-input v-model="routeHost"></b-form-input>
						</b-form-fieldset>
	
						<b-form-fieldset description="route prefix" label="Route Prefix" :state="routePrefixState">
							<b-form-input v-model="routePrefix"></b-form-input>
						</b-form-fieldset>
	
						<b-form-fieldset description="service name" label="Service Name" :state="serviceNameState">
							<b-form-input v-model="serviceName"></b-form-input>
						</b-form-fieldset>
	
						<b-form-fieldset description="service strip" label="Service Strip" state="success">
							<b-form-input v-model="serviceStrip"></b-form-input>
						</b-form-fieldset>
	
						<b-form-fieldset description="target host" label="Target Host" :state="targetHostState">
							<b-form-input v-model="targetHost"></b-form-input>
						</b-form-fieldset>
	
						<b-form-fieldset description="target port" label="Target Port" :state="targetPortState">
							<b-form-input v-model="targetPort"></b-form-input>
						</b-form-fieldset>
	
						<b-form-fieldset description="target weight" label="Target Weight" :state="targetWeightState">
							<b-form-input v-model="targetWeight"></b-form-input>
						</b-form-fieldset>
	
						<b-form-fieldset description="health check address" label="Health Check Address" :state="hcAddrState">
							<b-form-input v-model="hcAddr"></b-form-input>
						</b-form-fieldset>
						<b-form-fieldset description="health check interval" label="Health Check Interval" :state="hcIntervalState">
							<b-form-input v-model="hcInterval"></b-form-input>
						</b-form-fieldset>
					</div>
				</div>
				<!-- /.box-body -->
	
				<div class="box-footer">
					<b-button size="" variant="primary" v-on:click="submit">提交</b-button>
				</div>
			</div>
		</div>
	</section>
</template>

<script>
import { bFormFieldset, bFormInput, bButton } from 'bootstrap-vue/lib/components'
import axios from 'axios'
import config from '@/config'
import util from '@/util'

function state(bool) {
	return bool ? 'success' : 'warning'
}

export default {
	name: 'RouteForm',
	components: { bFormFieldset, bFormInput, bButton },
	data() {
		return {
			routeHost: '',
			routePrefix: '',
			serviceName: '',
			serviceStrip: '',
			targetHost: '',
			targetPort: 80,
			targetWeight: 10,
			hcAddr: '',
			hcInterval: '10s',
			alert:{
				class: 'alert-success',
				msg: 'OK',
				show: false
			}
		}
	},
	computed: {
		feedback() {
			return this.routeHost.length ? '' : 'emtpy';
		},
		routeHostState() {
			return state(this.routeHost);
		},
		routePrefixState() {
			return state(this.routePrefix);
		},
		serviceNameState() {
			return state(this.serviceName);
		},
		targetHostState() {
			return state(this.targetHost);
		},
		targetPortState() {
			return state(this.targetPort);
		},
		targetWeightState() {
			return state(this.targetWeight);
		},
		hcAddrState() {
			return state(this.hcAddr);
		},
		hcIntervalState() {
			return state(this.hcInterval);
		},
	},
	watch: {
		targetHost: function (val) {
			this.hcAddr = this.targetHost + ':' + this.targetPort
		},
		targetPort: function (val) {
			this.hcAddr = this.targetHost + ':' + this.targetPort
		}
	},
	methods: {
		submit: function () {
			var path = '/table/service';
			var domain = window.location.hostname;
        	var url = config.scheme + '://' + domain + ':' + config.port + path;
			var obj = this.$data
			var readyObj = util.camelToSnake(obj)
			var that = this

			readyObj.target_port = parseInt(readyObj.target_port)
			readyObj.target_weight = parseInt(readyObj.target_weight)
			console.log(readyObj)

			axios.put(url, readyObj, {
				headers: {'Content-Type': 'application/json'}
			}).then(function(response){
				console.log("then", response)
				if (response.data.code == 0) {
					that.alert.class = 'alert-success'
				} else {
					that.alert.class = 'alert-danger'
				}
			}).catch(e => {
				console.log("catching", e)
				if (e) {
					that.alert.class = 'alert-danger'
					that.alert.msg = e.response.data.msg
				}
			})

			this.alert.show = true
			console.log(this.alert)
			return false
		}
	}
}
</script>

<style>
#formWrapper {
	padding: 20px;
}
</style>