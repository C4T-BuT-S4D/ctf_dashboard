<template>
  <el-row :gutter="26">
    <el-col :span="8">
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
        <el-row :gutter="10">
          <el-col :span="12">
            <div class="text item">
              <a :href="boardLink" target="_blank">Scoreboard</a>
            </div>
            <div class="text item">
              <a :href="mongolLink" target="_blank">MonGol</a>
            </div>
            <div class="text item">
              <a :href="farmLink" target="_blank">Farm</a>
            </div>
            <div class="text item">
              Username:
              <span class="copiable" @click="copyText(username)">{{
                username
              }}</span>
            </div>
            <div class="text item">
              Password:
              <span class="copiable" @click="copyText(password)">{{
                password
              }}</span>
            </div>
          </el-col>
        </el-row>
      </el-card>
    </el-col>
    <el-col :span="8">
      <el-card>
        <div slot="header" class="clearfix">
          <span class="header-text">Files</span>
        </div>
        <el-row type="flex" justify="center">
          <el-col>
            <el-button
              v-for="file of files"
              :key="file"
              :span="6"
              class="item"
              @click="downloadFile(file)"
              >{{ file }}
            </el-button>
          </el-col>
        </el-row>
      </el-card>
    </el-col>
    <el-col :span="8">
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
            <a :href="getGoxyLink(vulnbox)" target="_blank">Goxy</a>
          </div>
          <ul>
            <li v-for="(service, j) of vulnbox.services" :key="j">
              {{ service.name }}:
              <a
                v-if="service.proto === 'http'"
                :href="getServiceLink(vulnbox, service)"
                target="_blank"
                >{{ getServiceLink(vulnbox, service) }}</a
              >
              <span
                v-else
                class="copiable"
                @click="copyText(getServiceLink(vulnbox, service))"
                >{{ getServiceLink(vulnbox, service) }}</span
              >
            </li>
          </ul>
        </div>
        <div>
          <el-row>
            <el-input
              v-model="keyContent"
              :rows="4"
              placeholder="Enter ssh public key"
              style="margin-bottom: 10px"
              type="textarea"
            ></el-input>
          </el-row>
          <el-row justify="center" type="flex">
            <el-button :span="8" @click="uploadKeyFile">Upload key</el-button>
          </el-row>
        </div>
      </el-card>
    </el-col>
  </el-row>
</template>

<script>
import moment from "moment";

export default {
  data: function() {
    return {
      config: null,
      username: "",
      password: "",
      vulnboxes: [],
      files: [],
      game: {},
      endTime: {},
      keyContent: "",
    };
  },

  mounted: async function() {
    await this.loadConfig();
    await this.loadFiles();
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
        this.endTime = moment(this.game.end);
        console.log(this.endTime);

        this.$notify({
          title: "Config load",
          message: "Success",
          type: "success",
          duration: 1000,
        });
      } catch (e) {
        this.config = {};
        this.$notify({
          title: "Config load",
          message: `Error: ${e}`,
          type: "error",
          duration: 10000,
        });
      }
    },

    loadFiles: async function() {
      try {
        const {
          data: { files },
        } = await this.$http.get("/files/");
        this.files = files;
      } catch (e) {
        this.$notify({
          title: "Loading files",
          message: `Error loading file list: ${e}`,
          type: "error",
          duration: 3000,
        });
      }
    },

    downloadFile: async function(name) {
      try {
        const response = await this.$http.get("/file/", {
          params: { name: name },
        });
        let fileURL = window.URL.createObjectURL(new Blob([response.data]));
        let fileLink = document.createElement("a");

        fileLink.href = fileURL;
        fileLink.setAttribute("download", name);
        document.body.appendChild(fileLink);

        fileLink.click();
      } catch (e) {
        this.$notify({
          title: "Download file",
          message: `Error downloading ${name}: ${e}`,
          type: "error",
          duration: 3000,
        });
      }
    },

    uploadKeyFile: async function() {
      try {
        await this.$http.post("/add_ssh_key/", { key: this.keyContent });
        this.$notify({
          title: "Key upload",
          message: "Success",
          type: "success",
        });
      } catch (e) {
        this.$notify({
          title: "Key upload",
          message: `Error: ${e}`,
          type: "error",
        });
      }
    },

    getGoxyLink: function(vulnbox) {
      return `http://${this.username}:${this.password}@${vulnbox.host}:${vulnbox.goxy_port}`;
    },

    getServiceLink: function(vulnbox, service) {
      if (service.proto === "http") {
        return `http://${vulnbox.host}:${service.port}`;
      } else {
        return `${vulnbox.host} ${service.port}`;
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
          duration: 1500,
        });
      } catch (err) {
        this.$notify({
          title: "Text copy failed",
          message: `Error: ${err}`,
          type: "error",
        });
      }

      document.body.removeChild(textArea);
    },
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
      return `http://${this.username}:${this.password}@${this.config.farm.addr}/?password=${this.password}`;
    },
    boardLink: function() {
      if (!this.config) {
        return "";
      }
      return this.game.board;
    },
  },
};
</script>

<style>
.text {
  font-size: 14px;
}

.item {
  margin-bottom: 18px !important;
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
