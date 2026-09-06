import service from '@/utils/request'

export function getModelConfigList(params) {
  return service({
    url: '/ai/modelConfig/list',
    method: 'GET',
    params
  })
}

export function createModelConfig(data) {
  return service({
    url: '/ai/modelConfig',
    method: 'POST',
    data
  })
}

export function updateModelConfig(data) {
  return service({
    url: '/ai/modelConfig',
    method: 'PUT',
    data
  })
}

export function deleteModelConfig(id) {
  return service({
    url: `/ai/modelConfig/${id}`,
    method: 'DELETE'
  })
}

export function setDefaultModelConfig(id) {
  return service({
    url: `/ai/modelConfig/setDefault/${id}`,
    method: 'PUT'
  })
}

export function getModelProviders() {
  return service({
    url: '/ai/modelConfig/providers',
    method: 'GET'
  })
}
