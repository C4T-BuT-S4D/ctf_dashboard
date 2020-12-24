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
        <div v-for="(vulnbox, i) of vulnboxes" :key="i">
          Addr: {{ vulnbox.host }}
          <a :href="getGoxyLink(vulnbox)">Go to goxy</a>
          <!-- <br /> -->
          <ul>
            <li v-for="(service_name, j) of vulnbox.services" :key="j">
              {{ service_name }}:
              <a :href="getServiceLink(vulnbox, service_name)"
                 v-if="service_map[service_name].proto === 'http'">{{ getServiceLink(vulnbox, service_name) }}</a>
              <span v-else class="copiable" @click="copyText(getServiceLink(vulnbox, service_name))">{{
                  getServiceLink(vulnbox, service_name)
                }}</span>
            </li>
          </ul>
        </div>
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
      vulnboxes: [],
      services: [],
      service_map: {}
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
        this.services = this.config.services;
        for (let service of this.services) {
          this.service_map[service.name] = service;
        }
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
      return `http://${this.username}:${this.password}@${vulnbox.host}:${vulnbox.goxy_port}`;
    },
    getServiceLink: function(vulnbox, service) {
      if (this.service_map[service].proto === "http") {
        return `http://${vulnbox.host}:${this.service_map[service].port}`;
      } else {
        return `${vulnbox.host} ${this.service_map[service].port}`;
      }
    },
    copyText: function(text) {
      var textArea = document.createElement("textarea");
      textArea.value = text;

      // Avoid scrolling to bottom
      textArea.style.top = "0";
      textArea.style.left = "0";
      textArea.style.position = "fixed";

      document.body.appendChild(textArea);
      textArea.focus();
      textArea.select();

      document.body.removeChild(textArea);

      this.$notify({title: "Text copied to clipboard", message: text, type: "success"});
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

.copiable {
  cursor: pointer;
  color: blue;
}
</style>
