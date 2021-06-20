import Vue from "vue";
import Element, {
  Button,
  Card,
  Col,
  Container,
  Header,
  Input,
  Main,
  Row
} from "element-ui";
import "element-ui/lib/theme-chalk/index.css";

Vue.use(Element);
Vue.use(Button);
Vue.component(Card.name, Card);
Vue.use(Container);
Vue.use(Main);
Vue.use(Header);
Vue.use(Row);
Vue.use(Col);
Vue.use(Input);
