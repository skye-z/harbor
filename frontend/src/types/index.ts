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
  repo_tags: string[]
  repo_digests: string[]
  created: number
  size: number
  shared_size: number
  virtual_size: number
  architecture: string
  os: string
}

export interface Network {
  id: string
  name: string
  driver: string
  scope: string
  attachable: boolean
  internal: boolean
  created: string
}

export interface Volume {
  id: string
  name: string
  driver: string
  mountpoint: string
  created: string
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
