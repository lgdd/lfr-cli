import App from './App.svelte'

const ELEMENT_ID = 'sample-custom-element-svelte';

if (!customElements.get(ELEMENT_ID)) {
  customElements.define(ELEMENT_ID, App.element);
}