import {createApp} from 'vue'
import App from './App.vue'

const ELEMENT_ID = 'sample-custom-element-vue';

class WebComponent extends HTMLElement {
  connectedCallback() {
    createApp(App).mount(this);
  }
}

if (!customElements.get(ELEMENT_ID)) {
  customElements.define(ELEMENT_ID, WebComponent);
}