// 音效管理工具

// 使用 Web Audio API 生成音效（无需外部音频文件）
class SoundManager {
  private audioContext: AudioContext | null = null

  private getAudioContext(): AudioContext {
    if (!this.audioContext) {
      this.audioContext = new AudioContext()
    }
    return this.audioContext
  }

  // 播放滴答声（滚动时）
  playTick(): void {
    try {
      const ctx = this.getAudioContext()
      const oscillator = ctx.createOscillator()
      const gainNode = ctx.createGain()

      oscillator.connect(gainNode)
      gainNode.connect(ctx.destination)

      oscillator.frequency.value = 800 + Math.random() * 400 // 随机频率 800-1200Hz
      oscillator.type = 'sine'

      gainNode.gain.setValueAtTime(0.1, ctx.currentTime)
      gainNode.gain.exponentialRampToValueAtTime(0.01, ctx.currentTime + 0.05)

      oscillator.start(ctx.currentTime)
      oscillator.stop(ctx.currentTime + 0.05)
    } catch (e) {
      console.warn('Failed to play tick sound:', e)
    }
  }

  // 播放结果揭晓音效
  playResult(): void {
    try {
      const ctx = this.getAudioContext()

      // 播放三个递进的音符
      const frequencies = [523.25, 659.25, 783.99] // C5, E5, G5

      frequencies.forEach((freq, i) => {
        const oscillator = ctx.createOscillator()
        const gainNode = ctx.createGain()

        oscillator.connect(gainNode)
        gainNode.connect(ctx.destination)

        oscillator.frequency.value = freq
        oscillator.type = 'sine'

        const startTime = ctx.currentTime + i * 0.1
        gainNode.gain.setValueAtTime(0.2, startTime)
        gainNode.gain.exponentialRampToValueAtTime(0.01, startTime + 0.3)

        oscillator.start(startTime)
        oscillator.stop(startTime + 0.3)
      })
    } catch (e) {
      console.warn('Failed to play result sound:', e)
    }
  }

  // 播放开始滚动音效
  playStart(): void {
    try {
      const ctx = this.getAudioContext()
      const oscillator = ctx.createOscillator()
      const gainNode = ctx.createGain()

      oscillator.connect(gainNode)
      gainNode.connect(ctx.destination)

      oscillator.frequency.setValueAtTime(300, ctx.currentTime)
      oscillator.frequency.exponentialRampToValueAtTime(600, ctx.currentTime + 0.2)
      oscillator.type = 'sawtooth'

      gainNode.gain.setValueAtTime(0.1, ctx.currentTime)
      gainNode.gain.exponentialRampToValueAtTime(0.01, ctx.currentTime + 0.2)

      oscillator.start(ctx.currentTime)
      oscillator.stop(ctx.currentTime + 0.2)
    } catch (e) {
      console.warn('Failed to play start sound:', e)
    }
  }
}

export const soundManager = new SoundManager()
