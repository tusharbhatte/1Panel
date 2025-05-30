<template>
    <div v-loading="loading">
        <LayoutContent title="mlBot 通知配置" :divider="true">
            <template #main>
                <el-row style="margin-top: 20px">
                    <el-col :span="1"><br /></el-col>
                    <el-col :xs="24" :sm="20" :md="20" :lg="12" :xl="12">
                        <el-form :model="form" label-position="left" ref="formRef" label-width="130px" :rules="rules">
                            <el-form-item label="启用通知" prop="enabled">
                                <el-switch v-model="form.enabled" @change="onEnabledChange" />
                                <span style="margin-left: 10px; color: #606266">启用后可发送 cronjob 执行结果通知</span>
                            </el-form-item>

                            <template v-if="form.enabled">
                                <el-form-item label="协议" prop="protocol">
                                    <el-select v-model="form.protocol" style="width: 200px">
                                        <el-option label="HTTP" value="http" />
                                        <el-option label="HTTPS" value="https" />
                                    </el-select>
                                </el-form-item>

                                <el-form-item label="主机地址" prop="host">
                                    <el-input
                                        v-model="form.host"
                                        placeholder="例如: 223.254.129.240"
                                        style="width: 300px"
                                    />
                                </el-form-item>

                                <el-form-item label="端口" prop="port">
                                    <el-input-number
                                        v-model="form.port"
                                        :min="1"
                                        :max="65535"
                                        placeholder="例如: 3000"
                                        style="width: 200px"
                                    />
                                </el-form-item>

                                <el-form-item label="用户ID" prop="uid">
                                    <el-input v-model="form.uid" placeholder="接收通知的用户ID" style="width: 300px" />
                                </el-form-item>

                                <el-form-item>
                                    <el-button type="primary" @click="onSave" :loading="saving">
                                        {{ $t('commons.button.save') }}
                                    </el-button>
                                    <el-button @click="onTest" :loading="testing" :disabled="!isFormValid">
                                        测试通知
                                    </el-button>
                                    <el-button type="info" @click="onInitialize" :loading="initializing">
                                        初始化配置
                                    </el-button>
                                </el-form-item>
                            </template>

                            <template v-else>
                                <el-form-item>
                                    <el-button type="primary" @click="onSave" :loading="saving">
                                        {{ $t('commons.button.save') }}
                                    </el-button>
                                </el-form-item>
                            </template>
                        </el-form>
                    </el-col>
                </el-row>
            </template>
        </LayoutContent>
    </div>
</template>

<script lang="ts" setup>
import { onMounted, reactive, ref, computed } from 'vue';
import { getMLBotConfig, updateMLBotConfig, testMLBotNotification } from '@/api/modules/mlbot';
import { MsgSuccess, MsgError } from '@/utils/message';
import i18n from '@/lang';

const loading = ref(false);
const saving = ref(false);
const testing = ref(false);
const initializing = ref(false);
const formRef = ref();

const form = reactive({
    enabled: false,
    host: '',
    port: 3000,
    uid: 'master',
    protocol: 'http',
});

const rules = {
    host: [{ required: true, message: '请输入主机地址', trigger: 'blur' }],
    port: [
        { required: true, message: '请输入端口', trigger: 'blur' },
        { type: 'number', min: 1, max: 65535, message: '端口范围: 1-65535', trigger: 'blur' },
    ],
    uid: [{ required: true, message: '请输入用户ID', trigger: 'blur' }],
};

const isFormValid = computed(() => {
    return form.enabled && form.host && form.port && form.uid;
});

const onEnabledChange = () => {
    if (!form.enabled) {
        // 清除验证错误
        formRef.value?.clearValidate();
    }
};

const onSave = async () => {
    if (form.enabled) {
        // 验证表单
        const valid = await formRef.value?.validate();
        if (!valid) return;
    }

    saving.value = true;
    try {
        await updateMLBotConfig({
            enabled: form.enabled,
            host: form.host,
            port: form.port,
            uid: form.uid,
            protocol: form.protocol,
        });
        MsgSuccess(i18n.global.t('commons.msg.operationSuccess'));
    } catch (error) {
        const errorMessage = error?.response?.data?.message || error?.message || error?.toString() || '未知错误';
        MsgError('保存失败: ' + errorMessage);
    } finally {
        saving.value = false;
    }
};

const onTest = async () => {
    // 验证表单
    const valid = await formRef.value?.validate();
    if (!valid) return;

    testing.value = true;
    try {
        await testMLBotNotification({
            host: form.host,
            port: form.port,
            uid: form.uid,
            protocol: form.protocol,
            message: '这是一条来自 1Panel 的测试通知消息',
        });
        MsgSuccess('测试通知发送成功');
    } catch (error) {
        const errorMessage = error?.response?.data?.message || error?.message || error?.toString() || '未知错误';
        MsgError('测试通知发送失败: ' + errorMessage);
    } finally {
        testing.value = false;
    }
};

const loadConfig = async () => {
    loading.value = true;
    try {
        const res = await getMLBotConfig();
        form.enabled = res.data.enabled;
        form.host = res.data.host;
        form.port = res.data.port;
        form.uid = res.data.uid;
        form.protocol = res.data.protocol;
    } catch (error) {
        const errorMessage = error?.response?.data?.message || error?.message || error?.toString() || '未知错误';
        MsgError('加载配置失败: ' + errorMessage);
    } finally {
        loading.value = false;
    }
};

const onInitialize = async () => {
    initializing.value = true;
    try {
        // 设置默认配置
        form.enabled = true;
        form.host = '223.254.129.240';
        form.port = 3000;
        form.uid = 'master';
        form.protocol = 'http';

        // 保存配置
        await updateMLBotConfig({
            enabled: form.enabled,
            host: form.host,
            port: form.port,
            uid: form.uid,
            protocol: form.protocol,
        });

        MsgSuccess('配置初始化成功');
    } catch (error) {
        const errorMessage = error?.response?.data?.message || error?.message || error?.toString() || '未知错误';
        MsgError('初始化失败: ' + errorMessage);
    } finally {
        initializing.value = false;
    }
};

onMounted(() => {
    loadConfig();
});
</script>
