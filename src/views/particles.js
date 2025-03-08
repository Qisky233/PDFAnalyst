const PARTICLE_COUNT = 60 // 粒子数量
const CONNECT_DISTANCE = 150 // 连接距离
const PARTICLE_SPEED = 0.3 // 粒子速度
const STAR_COLORS = ['#ffffff', '#f7f7a8', '#f0e68c', '#ffd700', '#ffa500'] // 星星颜色
const METEOR_RESET_DELAY = 3000 // 流星重置延迟（3秒）

let particles = []
let meteor = null // 只有一颗流星
let meteorTimeout = null // 流星重置的计时器
let animationFrameId = null

class Particle {
    constructor(canvas) {
        this.canvas = canvas
        this.ctx = canvas.getContext('2d')
        this.x = Math.random() * canvas.width
        this.y = Math.random() * canvas.height
        this.vx = (Math.random() - 0.5) * PARTICLE_SPEED
        this.vy = (Math.random() - 0.5) * PARTICLE_SPEED
        this.radius = Math.random() * 2 + 1 // 星星大小
        this.color = STAR_COLORS[Math.floor(Math.random() * STAR_COLORS.length)] // 随机星星颜色
        this.opacity = Math.random() * 0.8 + 0.2 // 随机透明度
    }

    update() {
        this.x += this.vx
        this.y += this.vy

        // 添加速度衰减
        this.vx *= 0.99
        this.vy *= 0.99

        // 边界检测，粒子超出边界后从另一边出现
        if (this.x < 0) this.x = this.canvas.width
        if (this.x > this.canvas.width) this.x = 0
        if (this.y < 0) this.y = this.canvas.height
        if (this.y > this.canvas.height) this.y = 0

        // 随机闪烁效果
        if (Math.random() < 0.01) {
            this.opacity = Math.random() * 0.8 + 0.2
        }
    }

    draw() {
        this.ctx.beginPath()
        this.ctx.arc(this.x, this.y, this.radius, 0, Math.PI * 2)
        this.ctx.fillStyle = `rgba(${parseInt(this.color.slice(1, 3), 16)}, ${parseInt(this.color.slice(3, 5), 16)}, ${parseInt(this.color.slice(5, 7), 16)}, ${this.opacity})`
        this.ctx.fill()
    }
}

class Meteor {
    constructor(canvas) {
        this.canvas = canvas
        this.ctx = canvas.getContext('2d')
        this.reset()
    }

    reset() {
        // 流星从右上角区域的随机位置出现
        const edge = Math.random() < 0.5 ? 'top' : 'right' // 50% 从顶部偏右出现，50% 从右侧偏上出现
        if (edge === 'top') {
            this.x = Math.random() * (this.canvas.width / 2) + (this.canvas.width / 2) // 上侧偏右
            this.y = -50 // 从顶部外开始
        } else {
            this.x = this.canvas.width + 50 // 从右侧外开始
            this.y = Math.random() * (this.canvas.height / 2) // 右侧偏上
        }

        // 固定的方向：从右上向左下
        this.vx = -Math.random() * 2 - 1 // 水平速度（向左，速度更慢）
        this.vy = Math.random() * 1 + 1 // 垂直速度（向下，速度更慢）

        // 随机尾巴长度
        this.length = Math.random() * 50 + 50 // 尾巴长度

        // 初始透明度
        this.opacity = 1

        // 标记为活跃状态
        this.active = true
    }

    update() {
        if (!this.active) return // 如果流星不活跃，则不更新

        this.x += this.vx
        this.y += this.vy

        // 流星飞出屏幕后重置
        if (
            this.y > this.canvas.height + this.length || // 飞出底部
            this.x < -this.length // 飞出左侧
        ) {
            this.active = false // 标记为非活跃状态
            this.resetWithDelay() // 延迟重置
        }

        // 流星逐渐消失
        this.opacity *= 0.98
    }

    resetWithDelay() {
        // 清除之前的计时器
        if (meteorTimeout) clearTimeout(meteorTimeout)

        // 3 秒后重置流星
        meteorTimeout = setTimeout(() => {
            this.reset()
        }, METEOR_RESET_DELAY)
    }

    draw() {
        if (!this.active) return // 如果流星不活跃，则不绘制

        // 绘制流星尾巴（朝后）
        const gradient = this.ctx.createLinearGradient(this.x, this.y, this.x - this.vx * 10, this.y - this.vy * 10)
        gradient.addColorStop(0, 'rgba(255, 255, 255, 0.8)')
        gradient.addColorStop(1, 'rgba(255, 255, 255, 0)')
        this.ctx.strokeStyle = gradient
        this.ctx.lineWidth = 2
        this.ctx.beginPath()
        this.ctx.moveTo(this.x, this.y)
        this.ctx.lineTo(this.x - this.vx * 10, this.y - this.vy * 10)
        this.ctx.stroke()
    }
}

export const initParticles = (canvas) => {
    // 确保 canvas 覆盖全屏
    canvas.width = window.innerWidth
    canvas.height = window.innerHeight

    // 初始化星星
    particles = Array(PARTICLE_COUNT).fill().map(() => new Particle(canvas))

    // 初始化一颗流星
    meteor = new Meteor(canvas)

    const animate = () => {
        const ctx = canvas.getContext('2d')
        ctx.clearRect(0, 0, canvas.width, canvas.height)

        // 绘制星星连接线
        particles.forEach(p1 => {
            particles.forEach(p2 => {
                const dx = p1.x - p2.x
                const dy = p1.y - p2.y
                if (dx * dx + dy * dy < CONNECT_DISTANCE * CONNECT_DISTANCE) {
                    ctx.beginPath()
                    ctx.moveTo(p1.x, p1.y)
                    ctx.lineTo(p2.x, p2.y)
                    ctx.strokeStyle = `rgba(255, 255, 255, ${0.2 - (dx * dx + dy * dy) / (CONNECT_DISTANCE * CONNECT_DISTANCE)})`
                    ctx.stroke()
                }
            })
        })

        // 更新并绘制星星
        particles.forEach(p => {
            p.update()
            p.draw()
        })

        // 更新并绘制流星
        if (meteor) {
            meteor.update()
            meteor.draw()
        }

        animationFrameId = requestAnimationFrame(animate)
    }

    animate()
}

export const destroyParticles = () => {
    cancelAnimationFrame(animationFrameId)
    if (meteorTimeout) clearTimeout(meteorTimeout)
}

// 窗口大小变化时重新初始化粒子
window.addEventListener('resize', () => {
    const canvas = document.querySelector('.particle-canvas')
    if (canvas) {
        canvas.width = window.innerWidth
        canvas.height = window.innerHeight
    }
})