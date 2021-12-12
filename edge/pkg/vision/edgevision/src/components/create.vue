<template>
    <div id="create">
        <!--导航条-->
        <el-menu :default-active="activeIndex2" class="el-menu-demo" mode="horizontal" @select="handleSelect"
                 background-color="#545c64" text-color="#fff" active-text-color="#ffd04b">
            <el-menu-item index="1">处理中心</el-menu-item>
            <el-submenu index="2">
                <template slot="title">我的工作台</template>
                <el-menu-item index="2-1">选项1</el-menu-item>
                <el-menu-item index="2-2">选项2</el-menu-item>
                <el-menu-item index="2-3">选项3</el-menu-item>
            </el-submenu>
            <el-menu-item index="3"><a href="https://www.ele.me" target="_blank">关于我们</a></el-menu-item>
        </el-menu>
        <br>
        <!--搜索框-->
        <el-row>
            <el-col :span="3" class="grid">
                <el-input v-model="input" placeholder="请输入内容" size="mini"></el-input>
            </el-col>
            <el-col :span="1" class="grid">
                <el-button type="success" icon="el-icon-search" size="mini">搜索</el-button>
            </el-col>
            <el-col :span="19" class="grid">
                <el-button type="success" round size="medium" style="float: right;" @click="dialogFormVisible = true">创建</el-button>
            </el-col>
        </el-row>
        <el-dialog title="创建" :visible.sync="dialogFormVisible">
            <el-form ref="form" :model="form" label-width="80px" size="medium">
                <el-form-item label="apiversion">
                  <el-input v-model="form.name"></el-input>
                </el-form-item>
                <el-form-item label="kind" label-width="50px">
                  <el-col :span="150" style="padding-left: 30px;">
                  <el-select v-model="form.region" placeholder="请选择创建类型" style="width: 100%;">
                    <el-option label="deployment" value="deployment"></el-option>
                    <el-option label="configmap" value="configmap"></el-option>
                  </el-select>
                </el-col>
                </el-form-item>
                <el-form-item label="metadata">
                  <el-col :span="11">
                    <el-input v-model="form.metadata.name" placeholder="Name"></el-input>
                  </el-col>
                  <el-col class="line" :span="2">-</el-col>
                  <el-col :span="11">
                    <el-input v-model="form.metadata.group" placeholder="Group"></el-input>
                  </el-col>
                </el-form-item>
                <el-form-item label="labels">
                  <el-col :span="11">
                    <el-input v-model="form.metadata.labels.app" placeholder="app"></el-input>
                  </el-col>
                  <el-col class="line" :span="2">-</el-col>
                  <el-col :span="11">
                    <el-input v-model="form.metadata.labels.rc" placeholder="rc"></el-input>
                  </el-col>
                </el-form-item>
                <el-form-item label="replicate">
                  <el-select v-model="form.spec.replicas" placeholder="请选择副本数量">
                    <el-option label="1" value="1"></el-option>
                    <el-option label="2" value="2"></el-option>
                  </el-select>
                </el-form-item>
                <el-form-item label="selector">
                  <el-col :span="11">
                    <el-input v-model="form.spec.selector.matchlabels.app" placeholder="app"></el-input>
                  </el-col>
                  <el-col class="line" :span="2">-</el-col>
                  <el-col :span="11">
                    <el-input v-model="form.spec.selector.matchlabels.rc" placeholder="rc"></el-input>
                  </el-col>
                </el-form-item>
                <el-form-item label="template">
                  <el-col :span="11">
                    <el-input v-model="form.spec.template.metadata.labels.app" placeholder="app"></el-input>
                  </el-col>
                  <el-col class="line" :span="2">-</el-col>
                  <el-col :span="11">
                    <el-input v-model="form.spec.template.metadata.labels.rc" placeholder="rc"></el-input>
                  </el-col>
                </el-form-item>
                <el-form-item label="spec">
                  <el-col :span="11">
                    <el-input v-model="form.spec.spec.containers.name" placeholder="name"></el-input>
                  </el-col>
                  <el-col class="line" :span="2">-</el-col>
                  <el-col :span="11">
                    <el-input v-model="form.spec.spec.containers.image" placeholder="imagepath"></el-input>
                  </el-col>
                </el-form-item>
                <el-form-item>
                  <el-button type="primary" @click="onSubmit">立即创建</el-button>
                  <el-button>取消</el-button>
                </el-form-item>
              </el-form>
            <!-- <div slot="footer" class="dialog-footer">
              <el-button @click="dialogFormVisible = false">取 消</el-button>
              <el-button type="primary" @click="dialogFormVisible = false">确 定</el-button>
            </div> -->

        </el-dialog>
        <br>
        <!--表格数据及操作-->
        <el-table :data="tableData" border style="width: 100%" stripe ref="multipleTable" tooltip-effect="dark">
            <!--勾选框-->
            <el-table-column type="selection" width="55">
            </el-table-column>
            <!--索引-->
            <el-table-column type="index" :index="indexMethod">
            </el-table-column>
            <el-table-column prop="date" label="日期" width="180" sortable>
            </el-table-column>
            <el-table-column prop="name" label="姓名" width="180">
            </el-table-column>
            <el-table-column prop="address" label="地址">
            </el-table-column>
            <el-table-column label="编辑" width="100">
                <template slot-scope="scope">
                    <el-button type="primary" icon="el-icon-edit" size="mini">编辑</el-button>
                </template>
            </el-table-column>
            <el-table-column label="删除" width="100">
                <template slot-scope="scope">
                    <el-button type="danger" icon="el-icon-delete" size="mini">删除</el-button>
                </template>
            </el-table-column>
        </el-table>
        <br>
        <!--新增按钮-->
        <el-col :span="1" class="grid">
            <el-button type="success" icon="el-icon-circle-plus-outline" size="mini" round>新增</el-button>
        </el-col>
        <!--全删按钮-->
        <el-col :span="1" class="grid">
            <el-button type="danger" icon="el-icon-delete" size="mini" round>全删</el-button>
        </el-col>
        <br>
        <!--分页条-->
        <el-pagination background layout="prev, pager, next" :total="1000">
        </el-pagination>
    </div>
</template>
<script>
import axios from 'axios'
    export default {
        data() {
            return {
                //表格数据
                tableData: [{
                    date: '2016-05-02',
                    name: '王小虎',
                    address: '上海市普陀区金沙江路 1518 弄'
                }, {
                    date: '2016-05-04',
                    name: '王小虎',
                    address: '上海市普陀区金沙江路 1517 弄'
                }, {
                    date: '2016-05-01',
                    name: '王小虎',
                    address: '上海市普陀区金沙江路 1519 弄'
                }, {
                    date: '2016-05-03',
                    name: '王小虎',
                    address: '上海市普陀区金沙江路 1516 弄'
                }],
                //查询输入框数据
                input: '',
                //导航条默认选项
                activeIndex: '1',
                activeIndex2: '1',
                gridData: [{
                date: '2016-05-02',
                name: '王小虎',
                address: '上海市普陀区金沙江路 1518 弄'
               }, {
          date: '2016-05-04',
          name: '王小虎',
          address: '上海市普陀区金沙江路 1518 弄'
        }, {
          date: '2016-05-01',
          name: '王小虎',
          address: '上海市普陀区金沙江路 1518 弄'
        }, {
          date: '2016-05-03',
          name: '王小虎',
          address: '上海市普陀区金沙江路 1518 弄'
        }],
        dialogFormVisible: false,
        form: {
          name: '',
          metadata:{
            name:'',
            labels:{
              app:'',
              rc:''
            }
          },
          spec:{
            replicas:'',
            selector:{
            matchlabels:{
              app:'',
              rc:''
            }
          },
          template:{
            metadata:{
              labels:{
                app:'',
                rc:''
              }
            }
          },
          spec:{
            containers:{
              name:'',
              image:''
            }
          },
          }
        },
        formLabelWidth: '120px',
        }
    },
        methods: {
            handleSelect(key, keyPath) {
                console.log(key, keyPath);
            },
            indexMethod(index) {
                return index;
            },
            onSubmit(){
              let formdata=JSON.stringify(this.form)
              var service = axios.create({
                baseURL:'/api',
                timeout:30000,
                headers:{
                  'content-type':'application/json'
                }
              });
             service.post('/create',formdata).then((success) => {
                console.log(success)
              }).catch((err) => {
                console.log(err)
              });
            }
        }
    }
</script>
<style>
    #create {
        font-family: Helvetica, sans-serif;
        text-align: center;
    }
</style>
