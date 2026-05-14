import Vue from 'vue'

/**
 * 防抖 单位时间只触发最后一次
 * 例：<el-button v-debounce="[reset,`click`,300]">刷新</el-button>
 * 简写：<el-button v-debounce="[reset]">刷新</el-button>
 */
Vue.directive('debounce', {
	inserted: function (el, binding) {
		let [fn, event = "click", time = 300] = binding.value
		let timer
		const handler = () => {
			timer && clearTimeout(timer)
			timer = setTimeout(() => fn(), time)
		}
		el.__debounceEvent__ = event
		el.__debounceHandler__ = handler
		el.__debounceClearTimer__ = () => {
			timer && clearTimeout(timer)
		}
		el.addEventListener(event, handler)
	},
	unbind: function (el) {
		if (el.__debounceHandler__) {
			el.removeEventListener(el.__debounceEvent__ || 'click', el.__debounceHandler__)
		}
		if (el.__debounceClearTimer__) {
			el.__debounceClearTimer__()
		}
		delete el.__debounceEvent__
		delete el.__debounceHandler__
		delete el.__debounceClearTimer__
	}
})

/**
 * 节流 每单位时间可触发一次
 * 例：<el-button v-throttle="[reset,`click`,300]">刷新</el-button>
 * 传递参数：<el-button v-throttle="[()=>reset(param),`click`,300]">刷新</el-button>
 */
Vue.directive('throttle', {
	inserted: function (el, binding) {
		let [fn, event = "click", time = 300] = binding.value
		let now, preTime
		const handler = () => {
			now = new Date()
			if (!preTime || now - preTime > time) {
				preTime = now
				fn()
			}
		}
		el.__throttleEvent__ = event
		el.__throttleHandler__ = handler
		el.addEventListener(event, handler)
	},
	unbind: function (el) {
		if (el.__throttleHandler__) {
			el.removeEventListener(el.__throttleEvent__ || 'click', el.__throttleHandler__)
		}
		delete el.__throttleEvent__
		delete el.__throttleHandler__
	}
})
