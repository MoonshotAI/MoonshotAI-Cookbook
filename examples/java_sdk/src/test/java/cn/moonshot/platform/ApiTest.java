package cn.moonshot.platform;

import cn.hutool.core.collection.CollUtil;
import cn.hutool.core.io.FileUtil;
import cn.moonshot.platform.util.Message;
import cn.moonshot.platform.util.MoonshotAiUtils;
import cn.moonshot.platform.util.RoleEnum;
import org.junit.jupiter.api.Test;

import java.util.List;


public class ApiTest {

    @Test
    void getModelList() {
        System.out.println(MoonshotAiUtils.getModelList());
    }

    @Test
    void uploadFile() {
        System.out.println(MoonshotAiUtils.uploadFile(FileUtil.file("/Users/steven/Desktop/test.pdf")));
    }

    @Test
    void getFileList() {
        System.out.println(MoonshotAiUtils.getFileList());
    }

    @Test
    void deleteFile() {
        System.out.println(MoonshotAiUtils.deleteFile("co17orilnl9coc91noh0"));
        System.out.println(MoonshotAiUtils.getFileList());
    }

    @Test
    void getFileContent() {
        System.out.println(MoonshotAiUtils.getFileContent("co18sokudu6bc6fqdhhg"));
    }

    @Test
    void getFileDetail() {
        System.out.println(MoonshotAiUtils.getFileDetail("co18sokudu6bc6fqdhhg"));
    }

    @Test
    void estimateTokenCount() {
        List<Message> messages = CollUtil.newArrayList(
                new Message(RoleEnum.system.name(), "你是kimi AI"),
                new Message(RoleEnum.user.name(), "hello")
        );
        System.out.println(MoonshotAiUtils.estimateTokenCount("moonshot-v1-8k", messages));
    }

    @Test
    void chat(){
        List<Message> messages = CollUtil.newArrayList(
                new Message(RoleEnum.system.name(), "你是kimi AI"),
                new Message(RoleEnum.user.name(), "hello")
        );
        MoonshotAiUtils.chat("moonshot-v1-8k",messages);
    }

}
