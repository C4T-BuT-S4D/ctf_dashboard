<template>
  <el-row :gutter="26">
    <el-col :span="6">
      <el-card>
        <div slot="header" class="clearfix">
          <span class="header-text">Game info</span>
          <countdown
            v-if="config"
            style="float: right; padding: 3px 0"
            :end-time="endTime"
          >
            <template v-slot:process="scope">
              <span>{{
                `${scope.timeObj.h}:${scope.timeObj.m}:${scope.timeObj.s}`
              }}</span>
            </template>
            <template v-slot:finish>
              <span>Game finished!</span>
            </template>
          </countdown>
        </div>
        <div class="row2">
          <div class="column2">
            <div class="text item"><a :href="boardLink">Scoreboard</a></div>
            <div class="text item">Username: {{ username }}</div>
            <div class="text item">Password: {{ password }}</div>
            <el-button @click="downloadKeyFile">Download key file</el-button>
          </div>
          <div class="column2">
            <div>
              <el-input
                :rows="4"
                type="textarea"
                placeholder="Enter ssh public key"
                v-model="keyContent"
              ></el-input>
              <el-button @click="uploadKeyFile">Upload key</el-button>
            </div>
          </div>
        </div>
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
          <div>
            <span
              class="copiable"
              @click="copyText(`${vulnbox.user}@${vulnbox.host}`)"
              >{{ vulnbox.user }}@{{ vulnbox.host }}</span
            >
            &rarr;
            <a :href="getGoxyLink(vulnbox)">Goxy</a>
          </div>
          <ul>
            <li v-for="(serviceName, j) of vulnbox.services" :key="j">
              {{ serviceName }}:
              <a
                :href="getServiceLink(vulnbox, serviceName)"
                v-if="serviceMap[serviceName].proto === 'http'"
                >{{ getServiceLink(vulnbox, serviceName) }}</a
              >
              <span
                v-else
                class="copiable"
                @click="copyText(getServiceLink(vulnbox, serviceName))"
                >{{ getServiceLink(vulnbox, serviceName) }}</span
              >
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
      vulnboxes: [],
      services: [],
      game: {},
      endTime: {},
      serviceMap: {},
      keyContent: ""
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
        this.vulnboxes = this.config.vulnboxes;

        this.game = this.config.game;
        this.endTime = new Date(this.game.end);
        console.log(this.endTime);

        this.services = this.config.services;
        for (let service of this.services) {
          this.serviceMap[service.name] = service;
        }
        this.$notify({
          title: "Config load",
          message: "Success",
          type: "success",
          duration: 1500
        });
      } catch (e) {
        this.config = {};
        this.$notify({
          title: "Config load",
          message: `Error: ${e}`,
          type: "error",
          duration: 10000
        });
      }
    },
    downloadFile: async function(path, name) {
      this.$http.get(path).then(response => {
        let fileURL = window.URL.createObjectURL(new Blob([response.data]));
        let fileLink = document.createElement("a");

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
    uploadKeyFile: async function() {
      try {
        await this.$http.post("/add_ssh_key/", { key: this.keyContent });
        this.$notify({
          title: "Key upload",
          message: "Success",
          type: "success"
        });
      } catch (e) {
        this.$notify({
          title: "Key upload",
          message: `Error: ${e}`,
          type: "error"
        });
      }
    },
    getGoxyLink: function(vulnbox) {
      return `http://${this.username}:${this.password}@${vulnbox.host}:${vulnbox.goxy_port}`;
    },
    getServiceLink: function(vulnbox, service) {
      if (this.serviceMap[service].proto === "http") {
        return `http://${vulnbox.host}:${this.serviceMap[service].port}`;
      } else {
        return `${vulnbox.host} ${this.serviceMap[service].port}`;
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

      try {
        document.execCommand("copy");
        this.$notify({
          title: "Text copied to clipboard",
          message: text,
          type: "success",
          duration: 1500
        });
      } catch (err) {
        this.$notify({
          title: "Text copy failed",
          message: `Error: ${err}`,
          type: "error"
        });
      }

      document.body.removeChild(textArea);
    }
  },
  computed: {
    mongolLink: function() {
      if (!this.config) {
        return "";
      }
      return `http://${this.username}:${this.password}@${this.config.mongol.addr}`;
    },
    farmLink: function() {
      if (!this.config) {
        return "";
      }
      return `http://${this.username}:${this.password}@${this.config.farm.addr}`;
    },
    boardLink: function() {
      if (!this.config) {
        return "";
      }
      return this.game.board;
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

.row2 {
  display: flex;
}

.column2 {
  flex: 50%;
}
</style>
