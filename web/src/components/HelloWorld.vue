
<template>
  <h1>奇淼的bilibili直播间弹幕</h1>
  <div v-for="(item,key) in MsgArr" :key="key">
    <span>{{item.userName}}</span>说:<span v-html="item.msg"></span>
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
  console.log(res)
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
    MsgArr.value.push({userName:resMsg.data.uname,msg:resMsg.data.msg})
    console.log(MsgArr.value)
  }
  
}

</script>

<style scoped>
a {
  color: #42b983;
}
</style>
