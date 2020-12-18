<template>
  <div id="app">
    <!-- <div id="nav">

    </div> -->

    <header class="grid-content header-color">
      <el-row>
        <a class="brand" href="#">CBS Dash</a>
        <!-- <router-link to="/">Home</router-link> |
        <router-link to="/about">About</router-link> -->
      </el-row>
    </header>
    <el-container>
      <el-main>
        <router-view />
      </el-main>
    </el-container>
  </div>
</template>

<style lang="scss">
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
}

#nav {
  padding: 30px;

  a {
    font-weight: bold;
    color: #2c3e50;

    &.router-link-exact-active {
      color: #42b983;
    }
  }
}

body {
  margin: 0px;
}

header {
  width: 100%;
  height: 60px;
  margin-bottom: 10px;
}

.header-color {
  background: #58b7ff;
}

#content {
  margin-top: 20px;
  padding-right: 40px;
}

.brand {
  color: #fff;
  background-color: transparent;
  margin-left: 20px;
  float: left;
  line-height: 25px;
  font-size: 25px;
  padding: 15px 15px;
  height: 30px;
  text-decoration: none;
}
</style>
