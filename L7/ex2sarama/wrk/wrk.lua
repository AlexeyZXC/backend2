math.randomseed(os.time())
request = function()
local k = math.random(0, 10)
local url = "/rate?rate="..k
return wrk.format("POST", url)
end
