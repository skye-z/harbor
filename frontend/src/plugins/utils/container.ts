export const containerStateMap: Record<string, string> = {
  running: '运行',
  exited: '停止',
  created: '创建',
  paused: '暂停',
  restarting: '重启',
  removing: '移除',
  dead: '异常'
}

export const formatContainerState = (state: string): string => {
  return containerStateMap[state] || state || '未知'
}

export const getContainerStateType = (state: string): 'success' | 'warning' | 'default' | 'error' => {
  if (state === 'running') return 'success'
  if (state === 'paused') return 'warning'
  return 'default'
}
