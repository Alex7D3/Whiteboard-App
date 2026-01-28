type RGB = [number,number,number];

interface LastWriteReg<T> {
  lastPeerID: number;
  timestamp: number;
  value: T|null;
  merge: (other: this) => void;
  set: (value: T) => void;
}

class PixelReg implements LastWriteReg<RGB> {
  readonly id: number;
  lastPeerID: number;
  timestamp: number;
  value: RGB|null;

  constructor({ id, lastPeerID, timestamp, value }: { id: number, lastPeerID: number, timestamp: number, value: RGB|null }) {
    this.id = id;
    this.lastPeerID = lastPeerID;
    this.timestamp = timestamp;
    this.value = value;
  }

  set(value: RGB|null) {
    this.timestamp++;
    this.lastPeerID = this.id;
    this.value = value;
  }

  merge(remote: LastWriteReg<RGB>) {
    const isRemoteNewer = this.timestamp < remote.timestamp;
    const tieBreak = this.timestamp === remote.timestamp && this.lastPeerID < remote.lastPeerID;

    if (isRemoteNewer || tieBreak) {
      this.timestamp = remote.timestamp;
      this.lastPeerID = remote.lastPeerID
      this.value = remote.value;
    }
  }


}

interface CrdtMap<T, R extends LastWriteReg<T>> {
  readonly id: number;
  data: Map<string, R>;
  merge: (remoteDiff: Map<string, R>) => void;
  set: (key: string, value: T) => void;
  get: (key: string) => T|undefined;
  has: (key: string) => boolean;
  delete: (key: string) => void;
}

class PixelCrdt implements CrdtMap<RGB, PixelReg> {
  readonly id: number;
  data: Map<string, PixelReg>;

  constructor(id: number) {
    this.id = id;
    this.data = new Map<string,PixelReg>();
  }

  set(key: string, value: RGB) {
    const reg = this.data.get(key);
    if (reg)
      reg.set(value);
    else
      this.data.set(key, new PixelReg({ id: this.id, lastPeerID: this.id, timestamp: 1, value }));
  }

  get(key: string): RGB {
    return this.data.get(key)?.value ?? [255,255,255];
  }

  has(key: string): boolean {
    return this.data.has(key);
  }

  delete(key: string) {
    this.data.get(key)?.set(null);
  }

  merge(remoteDiff: Map<string, PixelReg>) {
    for (const [key, remote] of remoteDiff.entries()) {
      const local = this.data.get(key);
      if (local)
        local.merge(remote);
      else
        this.data.set(key, new PixelReg({
          id: this.id,
          lastPeerID: remote.lastPeerID,
          timestamp: remote.timestamp,
          value: remote.value
        }));
    }
  }
}

export class CanvasManager {
  canvas: HTMLCanvasElement;
  ctx: CanvasRenderingContext2D;
  w: number;
  h: number;
  data: PixelCrdt;
  color: RGB;
  prevPoint: [number, number] | undefined;
  hasPainted: Set<string>;
  
  constructor(userID: number, canvas: HTMLCanvasElement) {
    this.canvas = canvas;
    const ctx = canvas.getContext("2d");
    if (!ctx) throw new Error("Failed to get 2d rendering context");
    this.ctx = ctx;
    this.w = canvas.width;
    this.h = canvas.height;
    this.data = new PixelCrdt(userID);
    this.color = [0,0,0];
    this.hasPainted = new Set<string>();

    this.canvas.addEventListener("pointerdown", this);
    this.canvas.addEventListener("pointermove", this);
    this.canvas.addEventListener("pointerup", this);

    this.canvas.width = this.canvas.clientWidth * devicePixelRatio;
    this.canvas.height = this.canvas.clientHeight * devicePixelRatio;
    this.ctx.scale(devicePixelRatio, devicePixelRatio);
    this.ctx.imageSmoothingEnabled = true;
    
  }

  key = (x: number, y: number) => `${x}${y}`;
  

  paint(x: number, y: number) {
    if (x < 0 || x >= this.w || y < 0 || y >= this.h)
      return;
    this.data.set(this.key(x, y), this.color);

    let [x0, y0] = this.prevPoint || [x, y];

    // DDA Algorithm for drawing a line
    const dx = x - x0, dy = y - y0;
    const steps = Math.max(Math.abs(dx), Math.abs(dy));
    const x_inc = dx / steps, y_inc = dy / steps;
    for (let _ = 0; _ < steps; _++) {
      x0 += x_inc;
      y0 += y_inc;
      const x1 = Math.round(x);
      const y1 = Math.round(y);
      this.data.set(this.key(x1, y1), this.color);
    }

    this.draw();
  }

  async draw() {
    const chans = 4;
    const buffer = new Uint8ClampedArray(this.w * this.h * chans);
    const row_size = this.w * chans;

    for (let y = 0; y < this.h; ++y) {
      const y_offset = y * row_size;
      for (let x = 0; x < this.w; ++x) {
        const x_offset = x * chans;

        const byte_offset = x_offset + y_offset;
        const [r, g, b] = this.data.get(`${x}${y}`);
        buffer[byte_offset] = r;
        buffer[byte_offset+1] = g;
        buffer[byte_offset+2] = b;
        buffer[byte_offset+3] = 255;
      }
    }
    const rawImage = new ImageData(buffer, this.w, this.h);
    const bitmap = await createImageBitmap(rawImage);
    this.ctx.drawImage(bitmap, 0, 0, this.w, this.h);
  }

  handleEvent(e: PointerEvent) {
    switch (e.type) {
      // @ts-expect-error
      case "pointerdown": {
        this.canvas.setPointerCapture(e.pointerId);
      }
      case "pointermove": {
        if (!this.canvas.hasPointerCapture(e.pointerId))
          return;

        const x = Math.floor((this.w * e.offsetX) / this.canvas.clientWidth); 
        const y = Math.floor((this.h * e.offsetY) / this.canvas.clientHeight);
        this.paint(x, y);
        this.prevPoint = [x, y];
        break;
      }

      case "pointerup": {
        this.canvas.releasePointerCapture(e.pointerId);
        this.hasPainted.clear();
        break;
      }
      
    }
  }
}
