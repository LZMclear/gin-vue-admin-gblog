import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import Fixture from './Fixture.vue'

createApp(Fixture).use(createPinia()).use(ElementPlus).mount('#app')
