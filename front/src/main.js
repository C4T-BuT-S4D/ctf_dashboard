import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";
import "@/plugins/element.js";
import "@/plugins/countdown.js";
import axios from "axios";
import { backUrl } from "@/config/index.js";

Vue.config.productionTip = false;
Vue.prototype.$http = axios;

axios.defaults.baseURL = backUrl;

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount("#app");
