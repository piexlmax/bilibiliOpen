
<template>
  <h1>奇淼的bilibili直播间弹幕</h1>
  <div v-for="(item,key) in MsgArr" :key="key" class="message">
    <img :src="item.uface" width="50" no-referrer>
    <span class="tag">{{item.fans_medal_name}}</span>
    <span class="name" :class="item.fans_medal_wearing_status?'vip':''">{{item.userName}}</span>
    说:<span v-html="item.msg"></span>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'

defineProps({
  msg: String
})
const count = ref(0)

const MsgArr = ref([])

const ws = new WebSocket("ws:127.0.0.1:8888/")

ws.onopen = (res)=>{
  const data = {
    cmd:"register",
    data:JSON.stringify({
      name:"test"+new Date()
    })
  }
  ws.send(JSON.stringify(data))
}

ws.onclose  = (res)=>{
  console.log(res,"连接关闭")
}

ws.onmessage = (res)=>{
  console.log(MsgArr.value)
  var reader = new FileReader();
  reader.readAsText(res.data, 'utf-8');
  reader.onload = function (e) {
    console.log(reader.result,"testtt")
    var resMsg = JSON.parse(reader.result)
    MsgArr.value.push({userName:resMsg.data.uname,msg:resMsg.data.msg,uface:resMsg.data.uface,fans_medal_name:resMsg.data.fans_medal_name,fans_medal_wearing_status:res.data.fans_medal_wearing_status})
    console.log(MsgArr.value)
  }
  
}

</script>

<style scoped>
a {
  color: #42b983;
}
.message{
  display: flex;
  align-items: center;
}
img{
  margin-right: 4px;
}

.tag{
  display: inline-flex;
  background: #42b983;
  padding: 2px 4px;
  font-size: 12px;
  margin-right: 4px;
}
.name{
  font-size: 14px;
}
.vip{
  color:red;
}
</style>
