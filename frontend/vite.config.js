import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

export default defineConfig({
    plugins: [vue()],
    resolve: {
        alias: {
            '@': fileURLToPath(new URL('./src', import.meta.url))
        }
    },
    // 生产环境使用相对路径（嵌入模式必须）
    base: './',
    build: {
        // 输出目录
        outDir: 'dist',
        // 静态资源目录
        assetsDir: 'assets',
        // 生成 sourcemap
        sourcemap: false,
        // 清空输出目录
        emptyOutDir: true
    },
    server: {
        port: 3000,
        proxy: {
            '/api': {
                target: 'http://localhost:8080',
                changeOrigin: true
            }
        }
    }
})
