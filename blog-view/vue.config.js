module.exports = {
	outputDir: 'dist',
	assetsDir: 'static',
	devServer: {
		proxy: {
			'/api': {
				target: 'http://127.0.0.1:8888/',
				changeOrigin: true,
				pathRewrite: {
					'^/api': ''
				}
			}
		}
	},
	configureWebpack: {
		resolve: {
			alias: {
				'assets': '@/assets',
				'common': '@/common',
				'components': '@/components',
				'api': '@/api',
				'views': '@/views',
				'plugins': '@/plugins'
			}
		}
	}
	
}
