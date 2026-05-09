module.exports = {
	outputDir: 'dist',
	assetsDir: 'static',
	devServer: {
		historyApiFallback: true,
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
	css: {
		loaderOptions: {
			css: {
				url: {
					filter: url => !url.startsWith('/') && !url.startsWith('data:')
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
