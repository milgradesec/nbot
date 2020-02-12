docker-machine.exe env rpi --shell powershell | Invoke-Expression 
docker stack deploy -c stack.yml nbot