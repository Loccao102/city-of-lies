// Procedural Web Audio API Sound Engine (Phase 11)
// Zero external files, zero latency, runs 100% in-browser offline

class SoundEngine {
  private ctx: AudioContext | null = null;
  private masterGain: GainNode | null = null;
  private tensionGain: GainNode | null = null;
  private tensionOsc1: OscillatorNode | null = null;
  private tensionOsc2: OscillatorNode | null = null;
  private tensionFilter: BiquadFilterNode | null = null;
  private isMuted: boolean = false;
  private isInitialized: boolean = false;
  private tensionStarted: boolean = false;

  private init() {
    if (typeof window === "undefined") return;
    if (this.isInitialized && this.ctx) {
      if (this.ctx.state === "suspended") {
        this.ctx.resume();
      }
      return;
    }

    try {
      const AudioCtx = window.AudioContext || (window as unknown as { webkitAudioContext: typeof AudioContext }).webkitAudioContext;
      this.ctx = new AudioCtx();
      this.masterGain = this.ctx.createGain();
      this.masterGain.gain.setValueAtTime(this.isMuted ? 0 : 0.4, this.ctx.currentTime);
      this.masterGain.connect(this.ctx.destination);
      this.isInitialized = true;
    } catch (e) {
      console.warn("Web Audio API not supported:", e);
    }
  }

  public resume() {
    this.init();
    if (this.ctx && this.ctx.state === "suspended") {
      this.ctx.resume();
    }
    if (!this.tensionStarted) {
      this.startAmbientTension();
    }
  }

  public toggleMute(): boolean {
    this.isMuted = !this.isMuted;
    if (this.masterGain && this.ctx) {
      this.masterGain.gain.setTargetAtTime(this.isMuted ? 0 : 0.4, this.ctx.currentTime, 0.05);
    }
    return this.isMuted;
  }

  public getIsMuted(): boolean {
    return this.isMuted;
  }

  public setVolume(vol: number) {
    if (this.masterGain && this.ctx && !this.isMuted) {
      const clamped = Math.max(0, Math.min(1, vol));
      this.masterGain.gain.setTargetAtTime(clamped, this.ctx.currentTime, 0.05);
    }
  }

  // Tactile button click
  public playClick() {
    this.init();
    if (!this.ctx || !this.masterGain || this.isMuted) return;

    const osc = this.ctx.createOscillator();
    const gain = this.ctx.createGain();

    osc.type = "sine";
    osc.frequency.setValueAtTime(900, this.ctx.currentTime);
    osc.frequency.exponentialRampToValueAtTime(300, this.ctx.currentTime + 0.04);

    gain.gain.setValueAtTime(0.2, this.ctx.currentTime);
    gain.gain.exponentialRampToValueAtTime(0.001, this.ctx.currentTime + 0.04);

    osc.connect(gain);
    gain.connect(this.masterGain);

    osc.start();
    osc.stop(this.ctx.currentTime + 0.04);
  }

  // Paper rustle sound for notebook tabs
  public playPaperRustle() {
    this.init();
    if (!this.ctx || !this.masterGain || this.isMuted) return;

    const bufferSize = this.ctx.sampleRate * 0.12; // 120ms
    const buffer = this.ctx.createBuffer(1, bufferSize, this.ctx.sampleRate);
    const data = buffer.getChannelData(0);
    for (let i = 0; i < bufferSize; i++) {
      data[i] = Math.random() * 2 - 1;
    }

    const noise = this.ctx.createBufferSource();
    noise.buffer = buffer;

    const filter = this.ctx.createBiquadFilter();
    filter.type = "bandpass";
    filter.frequency.setValueAtTime(1600, this.ctx.currentTime);
    filter.Q.setValueAtTime(3.0, this.ctx.currentTime);

    const gain = this.ctx.createGain();
    gain.gain.setValueAtTime(0.25, this.ctx.currentTime);
    gain.gain.exponentialRampToValueAtTime(0.001, this.ctx.currentTime + 0.12);

    noise.connect(filter);
    filter.connect(gain);
    gain.connect(this.masterGain);

    noise.start();
  }

  // Tech chime for rumor alert / public post
  public playRumorAlert() {
    this.init();
    if (!this.ctx || !this.masterGain || this.isMuted) return;

    const now = this.ctx.currentTime;
    const notes = [587.33, 880.0, 1174.66]; // D5, A5, D6 arpeggio

    notes.forEach((freq, idx) => {
      if (!this.ctx || !this.masterGain) return;
      const osc = this.ctx.createOscillator();
      const gain = this.ctx.createGain();

      osc.type = "triangle";
      osc.frequency.setValueAtTime(freq, now + idx * 0.06);

      gain.gain.setValueAtTime(0.18, now + idx * 0.06);
      gain.gain.exponentialRampToValueAtTime(0.001, now + idx * 0.06 + 0.25);

      osc.connect(gain);
      gain.connect(this.masterGain);

      osc.start(now + idx * 0.06);
      osc.stop(now + idx * 0.06 + 0.25);
    });
  }

  // Emergency vehicle siren (Fire engine / Ambulance)
  public playSiren() {
    this.init();
    if (!this.ctx || !this.masterGain || this.isMuted) return;

    const now = this.ctx.currentTime;
    const duration = 2.4; // 2.4 seconds

    const osc = this.ctx.createOscillator();
    const gain = this.ctx.createGain();
    const filter = this.ctx.createBiquadFilter();

    osc.type = "sawtooth";
    // Alternating two-tone siren 650Hz <-> 920Hz
    const cycles = 3;
    const period = duration / cycles;
    for (let c = 0; c < cycles; c++) {
      const t = now + c * period;
      osc.frequency.setValueAtTime(650, t);
      osc.frequency.linearRampToValueAtTime(920, t + period * 0.5);
      osc.frequency.linearRampToValueAtTime(650, t + period);
    }

    filter.type = "lowpass";
    filter.frequency.setValueAtTime(1400, now);

    gain.gain.setValueAtTime(0.01, now);
    gain.gain.linearRampToValueAtTime(0.15, now + 0.3);
    gain.gain.setValueAtTime(0.15, now + duration - 0.4);
    gain.gain.exponentialRampToValueAtTime(0.001, now + duration);

    osc.connect(filter);
    filter.connect(gain);
    gain.connect(this.masterGain);

    osc.start(now);
    osc.stop(now + duration);
  }

  // Fanfare on case solved / victory
  public playVictory() {
    this.init();
    if (!this.ctx || !this.masterGain || this.isMuted) return;

    const now = this.ctx.currentTime;
    const chords = [523.25, 659.25, 783.99, 1046.5]; // C5, E5, G5, C6
    chords.forEach((freq, idx) => {
      if (!this.ctx || !this.masterGain) return;
      const osc = this.ctx.createOscillator();
      const gain = this.ctx.createGain();

      osc.type = "sine";
      osc.frequency.setValueAtTime(freq, now + idx * 0.1);

      gain.gain.setValueAtTime(0.2, now + idx * 0.1);
      gain.gain.exponentialRampToValueAtTime(0.001, now + idx * 0.1 + 0.8);

      osc.connect(gain);
      gain.connect(this.masterGain);

      osc.start(now + idx * 0.1);
      osc.stop(now + idx * 0.1 + 0.8);
    });
  }

  // Somber gong on case defeat (false narrative reached 75%)
  public playDefeat() {
    this.init();
    if (!this.ctx || !this.masterGain || this.isMuted) return;

    const now = this.ctx.currentTime;
    const freqs = [110, 116.5, 130.8]; // Low dissonant cluster (A2, Bb2, C3)
    freqs.forEach((f) => {
      if (!this.ctx || !this.masterGain) return;
      const osc = this.ctx.createOscillator();
      const gain = this.ctx.createGain();

      osc.type = "sawtooth";
      osc.frequency.setValueAtTime(f, now);

      gain.gain.setValueAtTime(0.2, now);
      gain.gain.exponentialRampToValueAtTime(0.001, now + 2.5);

      osc.connect(gain);
      gain.connect(this.masterGain);

      osc.start(now);
      osc.stop(now + 2.5);
    });
  }

  // Dynamic Noir ambient tension drone
  private startAmbientTension() {
    if (!this.ctx || !this.masterGain || this.tensionStarted) return;
    this.tensionStarted = true;

    this.tensionFilter = this.ctx.createBiquadFilter();
    this.tensionFilter.type = "lowpass";
    this.tensionFilter.frequency.setValueAtTime(120, this.ctx.currentTime); // Deep dark bass rumble

    this.tensionGain = this.ctx.createGain();
    this.tensionGain.gain.setValueAtTime(0.12, this.ctx.currentTime);

    this.tensionOsc1 = this.ctx.createOscillator();
    this.tensionOsc1.type = "sine";
    this.tensionOsc1.frequency.setValueAtTime(55, this.ctx.currentTime); // A1 note (55 Hz)

    this.tensionOsc2 = this.ctx.createOscillator();
    this.tensionOsc2.type = "sine";
    this.tensionOsc2.frequency.setValueAtTime(56, this.ctx.currentTime); // 1Hz slow binaural beat

    this.tensionOsc1.connect(this.tensionFilter);
    this.tensionOsc2.connect(this.tensionFilter);
    this.tensionFilter.connect(this.tensionGain);
    this.tensionGain.connect(this.masterGain);

    this.tensionOsc1.start();
    this.tensionOsc2.start();
  }

  // Call whenever false narrative ratio changes
  public updateTension(ratio: number) {
    if (!this.ctx || !this.tensionFilter || !this.tensionGain) return;

    // As ratio climbs from 0.0 -> 0.75:
    // Lowpass filter opens from 120Hz up to 450Hz (more tense buzz)
    // Volume scales from 0.12 up to 0.28
    const clamped = Math.max(0, Math.min(0.8, ratio));
    const targetFreq = 120 + (clamped / 0.75) * 350;
    const targetGain = 0.12 + (clamped / 0.75) * 0.16;

    this.tensionFilter.frequency.setTargetAtTime(targetFreq, this.ctx.currentTime, 1.0);
    this.tensionGain.gain.setTargetAtTime(targetGain, this.ctx.currentTime, 1.0);
  }
}

export const sound = new SoundEngine();
