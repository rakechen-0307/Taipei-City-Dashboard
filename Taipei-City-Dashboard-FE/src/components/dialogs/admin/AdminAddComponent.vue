<!-- Developed by Taipei Urban Intelligence Center 2023-2024-->

<!-- This component only serves a functional purpose if a backend is connected -->
<!-- For static applications, this component could be removed or modified to be a dashboard component overviewer -->

<script setup>
import { computed, onMounted, ref } from "vue";
import http from "../../../router/axios";
import DashboardComponent from "../../../dashboardComponent/DashboardComponent.vue";

import { useDialogStore } from "../../../store/dialogStore";
import { useAdminStore } from "../../../store/adminStore";
import { useContentStore } from "../../../store/contentStore";

import DialogContainer from "../DialogContainer.vue";

const dialogStore = useDialogStore();
const adminStore = useAdminStore();
const contentStore = useContentStore();

const allComponents = ref(null);
const componentsSelected = ref([]);
const searchName = ref("");
const searchIndex = ref("");

// Filters out components already in the dashboard
const availableComponents = computed(() => {
	const taken = adminStore.currentDashboard.components?.map((item) => item.id) || [];
	const available = allComponents.value?.filter(
		(item) => !taken.includes(+item.id)
	);
	return available;
});

async function handleSearch() {
	const res = await http.get(`/component/`, {
		params: {
			pagesize: 200,
			searchbyindex: searchIndex.value,
			searchbyname: searchName.value,
			city: adminStore.currentCity,
		},
	});
	allComponents.value = res.data.data;
	contentStore.loading = false;
}
function handleSubmit() {
	adminStore.currentDashboard.components =
		adminStore.currentDashboard.components?.concat(componentsSelected.value) ?? componentsSelected.value;
	handleClose();
}
function handleClose() {
	searchName.value = "";
	searchIndex.value = "";
	componentsSelected.value = [];
	dialogStore.dialogs.adminAddComponent = false;
	handleSearch();
}

onMounted(() => {
	handleSearch();
});
</script>

<template>
  <DialogContainer
    :dialog="`adminAddComponent`"
    @on-close="handleClose"
  >
    <div class="addcomponent">
      <div class="addcomponent-header">
        <h2>新增組件至儀表板</h2>
        <div class="addcomponent-header-search">
          <div>
            <div>
              <input
                v-model="searchName"
                type="text"
                placeholder="以名稱搜尋 (Enter)"
                @keypress.enter="handleSearch"
              >
              <span
                v-if="searchName"
                @click="
                  () => {
                    searchName = '';
                    handleSearch();
                  }
                "
              >cancel</span>
            </div>
            <div>
              <input
                v-model="searchIndex"
                type="text"
                placeholder="以Index搜尋 (Enter)"
                @keypress.enter="handleSearch"
              >
              <span
                v-if="searchIndex"
                @click="
                  () => {
                    searchIndex = '';
                    handleSearch();
                  }
                "
              >cancel</span>
            </div>
          </div>
          <div>
            <button @click="handleClose">
              取消
            </button>
            <button
              v-if="componentsSelected?.length > 0"
              @click="handleSubmit"
            >
              <span>add_chart</span>確認新增
            </button>
          </div>
        </div>
      </div>
      <p :style="{ margin: '1rem 0 0.5rem' }">
        計 {{ availableComponents?.length }} 個組件符合篩選條件 | 共選取
        {{ componentsSelected?.length }} 個
      </p>

      <div class="addcomponent-list">
        <div
          v-for="item in availableComponents"
          :key="`${item.id}-${item.city}`"
        >
          <input
            :id="`${item.name}-${item.city}`"
            v-model="componentsSelected"
            type="checkbox"
            :value="{ id: item.id, name: item.name, city: item.city }"
          >
          <label :for="`${item.name}-${item.city}`">
            <div class="addcomponent-list-item">
              <DashboardComponent
                :config="item"
                :city-tag="contentStore.cityManager.getTagList(item.city)"
                mode="preview"
              />
            </div>
          </label>
        </div>
      </div>
    </div>
  </DialogContainer>
</template>

<style scoped lang="scss">
.addcomponent {
	width: 700px;
	height: 600px;
	padding: 10px;

	&-header {
		h2 {
			font-size: var(--font-m);
		}

		&-search {
			display: flex;
			justify-content: space-between;
			margin-top: var(--font-ms);

			> div {
				display: flex;
				justify-content: space-between;

				&:first-child {
					div {
						position: relative;
					}

					input {
						width: 150px;
						margin-right: 0.5rem;
					}

					span {
						position: absolute;
						right: 0.5rem;
						top: 0.4rem;
						margin-right: 4px;
						color: var(--color-complement-text);
						font-family: var(--font-icon);
						font-size: var(--font-m);
						transition: color 0.2s;
						cursor: pointer;

						&:hover {
							color: var(--color-highlight);
						}
					}
				}

				&:last-child {
					span {
						margin-right: 4px;
						font-family: var(--font-icon);
						font-size: calc(var(--font-m) * var(--font-to-icon));
					}

					button {
						display: flex;
						align-items: center;
						justify-self: baseline;
						margin-right: 0.4rem;
						border-radius: 5px;
						font-size: var(--font-m);

						&:nth-child(2) {
							padding: 2px 4px;
							background-color: var(--color-highlight);
						}
					}
				}
			}
		}
	}

	&-list {
		display: grid;
		grid-template-columns: 1fr 1fr;
		row-gap: var(--font-ms);
		column-gap: var(--font-ms);
		max-height: calc(100% - 7rem);
		overflow-y: scroll;

		&-item {
			border-radius: 5px;
			border: solid 1px var(--color-border);
			transition: border-color 0.2s, border-width 0.2s;
			cursor: pointer;
		}

		label {
			display: block;
		}

		input {
			display: none;
		}

		input:checked + label &-item {
			border-color: var(--color-highlight);
		}
	}
}
</style>
