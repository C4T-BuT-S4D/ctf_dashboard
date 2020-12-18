import Vue from "vue";
import Element from "element-ui";
import "element-ui/lib/theme-chalk/index.css";
import { Button, Card, Container, Main, Header, Row } from "element-ui";

Vue.use(Element);
Vue.use(Button);
Vue.component(Card.name, Card);
Vue.use(Container);
Vue.use(Main);
Vue.use(Header);
Vue.use(Row);
