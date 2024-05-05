<template>
  <v-container class="pa-0 ma-0 fill-height" fluid>
    <v-row>
      <v-col>
        <video
          v-if="videoUrl"
          :src="videoUrl"
          autoplay
          controls
          muted
          class="view"
        ></video>
        <v-img v-else :src="image" class="view"></v-img>
      </v-col>
    </v-row>
    <div
      class="position-absolute ma-1 font-weight-black text-white"
      :style="{ top: 0, right: 0, fontSize: '10px' }"
    >
      {{ isNewImage ? "NEUES FOTO" : currentImageCounter }}
    </div>
    <!-- <div
      class="position-absolute rounded-lg px-4 py-2 ma-1 font-weight-black"
      :style="{ background: '#fc5203', top: 0, right: 0 }"
    >
      {{ isNewImage ? "NEUES FOTO" : currentImageCounter }}
    </div> -->
  </v-container>
</template>

<script setup lang="ts">
const videoUrl = ref("");
const activeImage = ref(0);
const image = ref<string>();
const isNewImage = ref<boolean>();
const imageList = ref<string[]>([]);
const imageListBrandNew = ref<string[]>([]);
const imageTime = 5 * 1000; // show an image 5s
const newImageTime = 10 * 1000; // show new images 10s
const refreshTime = 3 * 1000; // update list every 3s

onMounted(async () => {
  // get all images and show the first one
  await listImages();

  // check if there is any saved number at local storage
  const counterString = localStorage.getItem("counter");
  if (counterString) {
    // start at saved position
    const counter = parseInt(counterString);
    activeImage.value = counter;
    await getImage(imageList.value[counter], false);
  } else {
    // start at 0
    await getImage(undefined, false);
  }

  // get list of images every few seconds -> there may are any new ones
  setInterval(() => {
    listImages();
  }, refreshTime);
});

const currentImageCounter = computed(() => {
  if (!imageList.value.length) return "0 / 0";
  return `${activeImage.value + 1} / ${imageList.value.length}`;
});

const listImages = async () => {
  // get all images
  const newImageList = await $fetch<string[]>(`/listImages`, {
    baseURL: "/api/v1",
    method: "GET",
  }).catch((err) => {
    console.error(err);
    return null;
  });
  if (newImageList !== null && newImageList.length !== imageList.value.length) {
    // we got images and the length changed -> update
    if (imageList.value.length) {
      // we got a list before -> so there are new images
      imageListBrandNew.value.push(
        ...newImageList.filter((item) => !imageList.value.includes(item))
      );
    }
    imageList.value = newImageList;
  }
};

const getImage = async (imageName: string | undefined, isNew: boolean) => {
  isNewImage.value = isNew;

  // check if an image exists
  if (!imageName) {
    // try to get new one
    imageName = imageListBrandNew.value.shift();
    if (imageName) {
      isNewImage.value = true;
    } else if (imageList.value.length) {
      // get first one
      imageName = imageList.value[0];
    } else {
      // probably no image exist right now -> try it again later
      setTimeout(() => getImage(undefined, false), imageTime);
      return;
    }
  }

  // fetch image
  const blob = await $fetch<Blob>(`/getImage/${imageName}`, {
    baseURL: "/api/v1",
    method: "GET",
  }).catch((err: any) => {
    console.error(err);
    return null;
  });
  if (blob) {
    if (blob.type.startsWith("video")) {
      videoUrl.value = URL.createObjectURL(blob);
    } else {
      videoUrl.value = "";
      const base64File = await blobToBase64(blob);
      image.value = base64File;
    }
  }

  setTimeout(
    () => {
      const newImage = imageListBrandNew.value.shift();
      if (newImage) {
        // get new added image
        getImage(newImage, true);
      } else {
        // prepare next regular image
        activeImage.value++;
        if (activeImage.value >= imageList.value.length) {
          activeImage.value = 0;
        }
        localStorage.setItem("counter", activeImage.value.toString());
        getImage(imageList.value[activeImage.value], false);
      }
    },
    isNewImage.value ? newImageTime : imageTime
  );
};

const blobToBase64 = async (blob: Blob) => {
  return new Promise<string>((resolve, _) => {
    const reader = new FileReader();
    reader.onloadend = () => resolve(reader.result as string);
    reader.readAsDataURL(blob);
  });
};
</script>

<style>
html {
  overflow-y: auto !important;
}
.view {
  min-width: 100%;
  height: calc(100vh - 32px);
}
</style>
