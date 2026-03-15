export interface Container {
  id: string
  names: string[]
  image: string
  image_id: string
  command: string
  created: string
  state: 'created' | 'running' | 'paused' | 'restarting' | 'removing' | 'exited' | 'dead'
  status: string
  ports: Port[]
  host_config: HostConfig
  network_settings: NetworkSettings
  size_rw?: number
  size_root_fs?: number
  mount?: Array<{ type: string; source: string; destination: string }>
  mounts?: Array<{ type: string; name?: string; source: string; destination: string }>
}

export interface Port {
  public_port?: number
  private_port: number
  type: string
}

export interface HostConfig {
  network_mode?: string
  port_bindings?: Record<string, Port[]>
  binds?: string[]
}

export interface NetworkSettings {
  networks: Record<string, NetworkInfo>
}

export interface NetworkInfo {
  network_id: string
  endpoint_id: string
  gateway: string
  ip_address: string
}

export interface Image {
  id: string
  parent_id: string
  repo_tags: string[]
  repo_digests: string[]
  created: number
  size: number
  virtual_size: number
  shared_size: number
  labels?: Record<string, string>
  containers: number
}

export interface Network {
  id: string
  name: string
  driver: string
  scope: string
  attachable: boolean
  internal: boolean
  created: string
  subnet?: string
  gateway?: string
  ipam?: {
    config?: Array<{
      subnet?: string
      gateway?: string
      ip_range?: string
    }>
  }
  labels?: Record<string, string>
}

export interface Volume {
  id: string
  name: string
  driver: string
  mountpoint: string
  created_at: string
  scope?: string
  usage_data?: {
    size: number
    ref_count: number
  }
  labels?: Record<string, string>
}

export interface DockerInfo {
  server_version: string
  operating_system: string
  os_type: string
  architecture: string
  ncpu: number
  mem_total: number
}

export interface ContainerStats {
  cpu_percent: number
  memory_usage: number
  memory_limit: number
  memory_percent: number
}

export interface UserInfo {
  id: string
  username: string
  is_admin: boolean
}

export interface WSMessage {
  type: string
  payload: any
}
