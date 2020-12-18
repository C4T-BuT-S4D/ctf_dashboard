<template>
  <el-row :gutter="26">
    <el-col :span="6">
      <el-card>
        <div slot="header" class="clearfix">
          <span class="header-text">Credentials</span>
        </div>
        <div class="text item">Username: {{ username }}</div>
        <div class="text item">Password: {{ password }}</div>
        <el-button @click="downloadKeyFile">Download key file</el-button>
      </el-card>
    </el-col>
    <el-col :span="6">
      <el-card>
        <div slot="header" class="clearfix">
          <span class="header-text">Mongol</span>
          <a :href="mongolLink" style="float: right; padding: 3px 0" type="text"
            >Open</a
          >
        </div>
        <div class="text item">Credentials are the same</div>
      </el-card>
    </el-col>
    <el-col :span="6">
      <el-card>
        <div slot="header" class="clearfix">
          <span class="header-text">Farm</span>
          <a :href="farmLink" style="float: right; padding: 3px 0" type="text"
            >Open</a
          >
        </div>
        <el-button @click="downloadStartSploit"
          >Download start sploit</el-button
        >
      </el-card>
    </el-col>
    <el-col :span="6">
      <el-card>
        <div slot="header" class="clearfix">
          <span class="header-text">Vulnboxes</span>
        </div>
        <p v-for="(vulnbox, i) of vulnboxes" :key="i">
          Addr: {{ vulnbox.host }}
          <a :href="getGoxyLink(vulnbox)">Go to goxy</a>
        </p>
      </el-card>
    </el-col>
  </el-row>
</template>

<script>
export default {
  data: function() {
    return {
      config: null,
      username: "",
      password: "",
      mongolLink: "",
      farmLink: "",
      vulnboxes: []
    };
  },

  mounted: async function() {
    await this.loadConfig();
  },

  methods: {
    loadConfig: async function() {
      try {
        const { data } = await this.$http.get("/config/");
        this.config = data;
        this.username = this.config.auth.username;
        this.password = this.config.auth.password;
        this.mongolLink = `http://${this.username}:${this.password}@${this.config.mongol.addr}`;
        this.farmLink = `http://${this.username}:${this.password}@${this.config.farm.addr}`;
        this.vulnboxes = this.config.vulnboxes;
      } catch {
        this.config = {};
      }
    },
    downloadFile: async function(path, name) {
      this.$http.get(path).then(response => {
        var fileURL = window.URL.createObjectURL(new Blob([response.data]));
        var fileLink = document.createElement("a");

        fileLink.href = fileURL;
        fileLink.setAttribute("download", name);
        document.body.appendChild(fileLink);

        fileLink.click();
      });
    },
    downloadKeyFile: async function() {
      await this.downloadFile("/key_file", "ssh_key");
    },
    downloadStartSploit: async function() {
      await this.downloadFile("/start_sploit.py", "start_sploit.py");
    },
    getGoxyLink: function(vulnbox) {
      console.log(vulnbox);
      return `http://${this.username}:${this.password}@${vulnbox.host}:${vulnbox.goxy_port}`;
    }
  }
};
</script>

<style>
.text {
  font-size: 14px;
}

.item {
  margin-bottom: 18px;
}

.header-text {
  font-weight: 900;
}

.clearfix:before,
.clearfix:after {
  display: table;
  content: "";
}
.clearfix:after {
  clear: both;
}
</style>
