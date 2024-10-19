<template>
  <div class="todo-main">
    <h1>TODOリスト</h1>

    <!-- Show status messages -->
    <div v-if="statusMessage" class="status-message">{{ statusMessage }}</div>

    <!-- Input group to add tasks -->
    <div class="input-group">
      <input
        v-model="newTask"
        placeholder="新しいタスクを入力 (Enter new task)"
        @keyup.enter="addTodo"
        class="left-panel"
      />
      <select v-model="newPriority" class="priority-dropdown right-panel" >
        <option value="high">高 (High)</option>
        <option value="medium">中 (Medium)</option>
        <option value="low">低 (Low)</option>
      </select>
      <button class="green-button" @click="addTodo">追加 (Add)</button>
    </div>

    <!-- Search Section -->
    <div class="search-section">
      <select v-model="searchType" class="search-dropdown same-panel" @change="resetSearchInputs">
        <option value="task">タスク (Task)</option>
        <option value="status">ステータス (Status)</option>
      </select>
      <input
        v-if="searchType === 'task'"
        v-model="searchTask"
        placeholder="タスクを検索 (Search Task)"
        class="search-input same-panel"
        @input="filterTodos"
      />
      <select
        v-if="searchType === 'status'"
        v-model="searchStatus"
        class="search-dropdown"
        @change="filterTodos"
      >
        <option value="">すべて (All)</option>
        <option value="created">作成済み (Created)</option>
        <option value="processing">処理中 (Processing)</option>
        <option value="done">完了 (Done)</option>
      </select>
      <button class="green-button" @click="searchType === 'task' ? searchByTask() : searchByStatus()">
        検索 (Search)
      </button>
    </div>

    <!-- Display tasks -->
    <div v-if="filteredTodos.length > 0">
      <div v-for="todo in filteredTodos" :key="todo.ID" class="todo-item">
        <div>
          <span
            v-if="!todo.isEditing"
            :class="{ 'done-task': todo.Status === 'done' }"
          >
            {{ todo.Task }} (Status : {{ todo.Status }}) (Priority : {{ todo.Priority }})
          </span>
          <div v-else class="edit-dropdown">
            <input
              v-model="todo.Task"
              class="edit-input"
              placeholder="Edit Task"
            />
            <select
              v-model="todo.Status"
              class="status-dropdown"
            >
              <option value="created">作成済み (Created)</option>
              <option value="processing">処理中 (Processing)</option>
              <option value="done">完了 (Done)</option>
            </select>
            <button class="submit-button" @click="updateTodo(todo)">更新 (Update)</button> 
          </div>
        </div>
        <div class="buttons">
          <button class="edit-button" @click="enableEdit(todo)">✏️</button>
          <button class="delete-button" @click="deleteTodo(todo.ID)">🗑️</button>
        </div>
      </div>
    </div>
    <div v-else>
      <p>タスクがありません (No tasks available).</p>
    </div>
  </div>
</template>

<script lang="ts">
import { fetchTodos, createTodo, updateTodo, deleteTodo } from '../services/TodoService.js';

interface Todo {
  ID: number;
  Task: string;
  Status: string;
  Priority: string;
  isEditing?: boolean;
}

export default {
  data() {
    return {
      newTask: '',
      newPriority: 'low',
      todos: [] as Todo[],
      statusMessage: '',
      searchType: 'task',
      searchTask: '',
      searchStatus: '',
      filteredTodos: [] as Todo[],
    };
  },
  mounted() {
    this.loadTodos();
  },
  methods: {
    async loadTodos() {
      try {
        const data = await fetchTodos();
        this.todos = data;
        this.filteredTodos = this.todos;
      } catch (error) {
        this.setStatusMessage('タスクの取得に失敗しました (Failed to fetch tasks)');
      }
    },
    async addTodo() {
      if (!this.newTask.trim()) return;
      try {
        await createTodo({ task: this.newTask, Status: 'created', Priority: this.newPriority });
        this.newTask = '';
        this.newPriority = 'low';
        await this.loadTodos();
        this.setStatusMessage('タスクが追加されました (Task added)');
      } catch (error) {
        this.setStatusMessage('タスクの作成に失敗しました (Failed to create task)');
      }
    },
    async editTodo(todo: Todo) {
      todo.isEditing = false;
      try {
        await updateTodo(todo.ID, todo);
        this.setStatusMessage('タスクが編集されました (Task edited)');
        await this.loadTodos();
      } catch (error) {
        this.setStatusMessage('タスクの編集に失敗しました (Failed to edit task)');
      }
    },
    enableEdit(todo: Todo) {
      // Set all tasks to non-editing mode
      this.todos.forEach(t => {
        t.isEditing = false;
      });
      
      // Set the selected task to editing mode
      todo.isEditing = true;
    },
    async deleteTodo(id: number) {
      console.log("Attempting to delete todo with ID:", id); 
      try {
        await deleteTodo(id); // Attempt to delete the task
        await this.loadTodos(); // Refresh the list after deletion
        this.setStatusMessage('タスクが削除されました (Task deleted)'); 
      } catch (error) {
        console.error("Error deleting task:", error); 
        this.setStatusMessage('タスクの削除に失敗しました (Failed to delete task)'); 
      }
    },

    async updateTodo(todo: Todo) {
      // Exit editing mode
      todo.isEditing = false; 
      
      try {
        // Update the task and its status via the API
        await updateTodo(todo.ID, {
          Task: todo.Task,
          Status: todo.Status // Include the status in the update
        });
        this.setStatusMessage('タスクが更新されました (Task updated)');
        await this.loadTodos(); // Reload the tasks
      } catch (error) {
        this.setStatusMessage('タスクの更新に失敗しました (Failed to update task)');
      }
    },
    filterTodos() {
      this.filteredTodos = this.todos.filter((todo: Todo) => {
        const matchesTask = todo.Task.toLowerCase().includes(this.searchTask.toLowerCase());
        const matchesStatus = this.searchStatus ? todo.Status === this.searchStatus : true;
        return matchesTask && matchesStatus;
      });
    },
    async searchByTask() {
      this.filterTodos(); // Call the filter method
    },
    async searchByStatus() {
      this.filterTodos(); // Call the filter method
    },
    resetSearchInputs() {
      this.searchTask = '';
      this.searchStatus = '';
      this.filteredTodos = this.todos;
    },
    setStatusMessage(message: string) {
      this.statusMessage = message;
      setTimeout(() => {
        this.statusMessage = '';
      }, 5000);
    },
  },
};
</script>

<style scoped>
@import '../assets/styles/TodoList.css';
</style>
