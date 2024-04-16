Rails.application.routes.draw do
  get '/hello', to: 'hello#local_hello'

  get '/remotehello', to: 'hello#remote_hello'
end
