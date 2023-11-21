import httpClient from '../utils/httpClient'

export default {
  list (query, callback) {
    httpClient.get('/task/log', query, callback)
  },

  clear (callback) {
    httpClient.post('/task/log/clear', {}, callback)
  },

  stop (id, taskId, callback) {
    httpClient.post('/task/log/stop', {id, task_id: taskId}, callback)
  },

  remove(month,callback){
    httpClient.post(`/task/log/remove/${month}`,{},callback)
  },

  removeDay(day,callback){
    httpClient.post(`/task/log/remove/day/${day}`,{},callback)
  }
}
