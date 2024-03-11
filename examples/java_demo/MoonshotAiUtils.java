package moonshot.example;

import com.alibaba.fastjson.JSON;
import com.alibaba.fastjson.JSONObject;
import com.dh.bigdata.invest.tools.utils.HttpClientUtil;
import com.dh.bigdata.invest.tools.utils.HttpsClientUtil;
import com.dh.bigdata.invest.tools.utils.JacksonUtils;
import org.apache.commons.lang3.StringUtils;
import org.springframework.web.multipart.MultipartFile;

import java.io.IOException;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class MoonshotAiUtils {
	public static String BASE_COMPLETION_URL = "https://api.moonshot.cn/v1";
	public static String CHAT_COMPLETION_URL = "https://api.moonshot.cn/v1/chat/completions";
	public static String CHAT_LIST_MODELS_URL = "https://api.moonshot.cn/v1/models";
	public static String UPLOAD_FILES_URL = "https://api.moonshot.cn/v1/files";
	public static String ESTIMATE_TOKEN_COUNT_URL = "https://api.moonshot.cn/v1/tokenizers/estimate-token-count";
	public static String MOONSHOT_API_KEY = "";
	public static String API_MODEL_8K = "moonshot-v1-8k";
	public static String SYSTEM_CONTENT = "你是 Kimi，由 Moonshot AI 提供的人工智能助手，你更擅长中文和英文的对话。你会为用户提供安全，有帮助，准确的回答。同时，你会拒绝一些涉及恐怖主义，种族歧视，黄色暴力等问题的回答。Moonshot AI 为专有名词，不可翻译成其他语言。";

	public static void main(String[] args) {
		JSONObject jsonObject = httpOutputToPage( "123",  "请参考刚上传的文件",  "cnjc0k6cp7f2ounekiag",  "");
	}

	public static void sendChatCompletions(String question) {
		sendChatCompletionsCommon(question,CHAT_COMPLETION_URL);
	}

	public static void sendChatCompletionsTokenCount(String question) {
		sendChatCompletionsCommon(question, ESTIMATE_TOKEN_COUNT_URL);
	}

	public static void sendChatCompletionsCommon(String question,String url) {
		Map<String, String> headerMap = new HashMap<>();
		headerMap.put("Content-Type", "application/json");
		headerMap.put("Authorization", "Bearer " + MOONSHOT_API_KEY);
		Map<String, Object> paramMap = new HashMap<>();
		paramMap.put("model", API_MODEL_8K);
		paramMap.put("temperature", 0.2);
		paramMap.put("max_tokens", 4096);
		List<Map<String, String>> messageList = new ArrayList<>();
		Map<String, String> sysMsgMap = new HashMap<>();
		sysMsgMap.put("role", "system");
		Map<String, String> userMsgMap = new HashMap<>();
		userMsgMap.put("role", "user");
		String sysContent = SYSTEM_CONTENT;
		String userContent = "推荐100家企业清单";
		if (StringUtils.isNotBlank(question)) {
			userContent = question;
		}
		sysMsgMap.put("content", sysContent);
		userMsgMap.put("content", userContent);
		messageList.add(sysMsgMap);
		messageList.add(userMsgMap);
		paramMap.put("messages", messageList);
//		paramMap.put("stream", true);

		String result = null;
		try {
			result = HttpClientUtil.postAddHeadersWithObjEntity(url, paramMap, headerMap);

			System.out.println(result);
		} catch (IOException e) {
			e.printStackTrace();
		}
	}

	/**
	 * @Description list models-模型列表
	 * 你可以使用我们的 List Models API 来获取当前可用的模型列表。
	 *
	 * 当前的，我们支持的模型有：
	 *
	 * moonshot-v1-8k: 它是一个长度为 8k 的模型，适用于生成短文本。
	 * moonshot-v1-32k: 它是一个长度为 32k 的模型，适用于生成长文本。
	 * moonshot-v1-128k: 它是一个长度为 128k 的模型，适用于生成超长文本。
	 * 以上模型的区别在于它们的最大上下文长度，这个长度包括了输入消息和生成的输出，在效果上并没有什么区别。这个主要是为了方便用户选择合适的模型。
	 * @param  ->
	 * @return
	 **/
	public static Map<String, Object> listModels(){
		Map<String, String> headerMap = new HashMap<>();
		headerMap.put("Authorization", "Bearer " + MOONSHOT_API_KEY);
		String result = null;
		try {
			result = HttpsClientUtil.getAddHeaders(CHAT_LIST_MODELS_URL, headerMap);

			System.out.println(result);
		} catch (IOException e) {
			e.printStackTrace();
		}
		return JacksonUtils.fromJsonToMap(result);
	}
	/**
	 * @param ->
	 * @return
	 * @Description 上传文件
	 **/
	public static Map<String, Object> uploadMSAiFiles(MultipartFile file) {
		Map<String, String> headerMap = new HashMap<>();
//		headerMap.put("Content-Type", "multipart/form-data");
		headerMap.put("Authorization", "Bearer " + MOONSHOT_API_KEY);
		Map<String, String> paramMap = new HashMap<>();
		paramMap.put("purpose","file-extract");
		String result = null;
		try {
			result = HttpsClientUtil.postAddHeadersFiles(UPLOAD_FILES_URL, paramMap, headerMap, file);

			System.out.println(result);
		} catch (IOException e) {
			e.printStackTrace();
		}
		return JacksonUtils.fromJsonToMap(result);
	}

	/**
	 * @Description 文件列表
	 * @param  ->
	 * @return
	 **/
	public static Map<String, Object> filesList(){
		Map<String, String> headerMap = new HashMap<>();
		headerMap.put("Authorization", "Bearer " + MOONSHOT_API_KEY);
		String result = null;
		try {
			result = HttpsClientUtil.getAddHeaders(UPLOAD_FILES_URL, headerMap);

			System.out.println(result);
		} catch (IOException e) {
			e.printStackTrace();
		}
		return JacksonUtils.fromJsonToMap(result);
	}

	/**
	 * @Description 根据API 文件id查看单个文件
	 * @param  ->
	 * @return
	 **/
	public static Map<String, Object> fileInfo(String fileId){
		Map<String, String> headerMap = new HashMap<>();
		headerMap.put("Authorization", "Bearer " + MOONSHOT_API_KEY);
		String result = null;
		try {
			result = HttpsClientUtil.getAddHeaders(UPLOAD_FILES_URL+"/"+fileId, headerMap);

			System.out.println(result);
		} catch (IOException e) {
			e.printStackTrace();
		}
		return JacksonUtils.fromJsonToMap(result);
	}

	/**
	 * @Description 根据文件id查看文件内容
	 * @param  ->
	 * @return
	 **/
	public static Map<String, Object> fileContent(String fileId){
		Map<String, String> headerMap = new HashMap<>();
		headerMap.put("Authorization", "Bearer " + MOONSHOT_API_KEY);
		String result = null;
		try {
			result = HttpsClientUtil.getAddHeaders(UPLOAD_FILES_URL+"/"+fileId+"/content", headerMap);

			System.out.println(result);
		} catch (IOException e) {
			e.printStackTrace();
		}
		return JacksonUtils.fromJsonToMap(result);
	}

	/**
	 * @Description 根据文件id删除文件
	 * @param  ->
	 * @return
	 **/
	public static Map<String, Object> delfile(String fileId){
		Map<String, String> headerMap = new HashMap<>();
		headerMap.put("Authorization", "Bearer " + MOONSHOT_API_KEY);
		String result = null;
		try {
			result = HttpsClientUtil.deleteAddHeaders(UPLOAD_FILES_URL+"/"+fileId, headerMap);

			System.out.println(result);
		} catch (IOException e) {
			e.printStackTrace();
		}
		return JacksonUtils.fromJsonToMap(result);
	}

	static HashMap<String, String> assistantMap = new HashMap<>();

	public static JSONObject httpOutputToPage(String clientId, String question, String fileId, String goOnContent) {
		JSONObject jsonObject = null;
		Map<String, String> headerMap = new HashMap<>();
		headerMap.put("Content-Type", "application/json");
		headerMap.put("Authorization", "Bearer " + MOONSHOT_API_KEY);
		Map<String, Object> paramMap = new HashMap<>();
		paramMap.put("model", API_MODEL_8K);
		paramMap.put("temperature", 0.2);
		paramMap.put("max_tokens", 4096);
		List<Map<String, String>> messageList = new ArrayList<>();
		Map<String, String> sysMsgMap = new HashMap<>();
		sysMsgMap.put("role", "system");
		Map<String, String> userMsgMap = new HashMap<>();
		userMsgMap.put("role", "user");
		String sysContent = SYSTEM_CONTENT;
		String userContent = "园区招商推荐100家企业清单";
		sysMsgMap.put("content", sysContent);
		if (StringUtils.isNotBlank(question)) {
			userContent = question;
		}
		userMsgMap.put("content", userContent);
		messageList.add(sysMsgMap);
		messageList.add(userMsgMap);
		paramMap.put("messages", messageList);
		String requestUrl = CHAT_COMPLETION_URL;
		if(StringUtils.isNotBlank(fileId)){
			Map<String, Object> fileContent=fileContent(fileId);
			if(fileContent!=null&&fileContent.size()>0&&StringUtils.isNotBlank((String)fileContent.get("content"))){
				Map<String, String> fileMsgMap = new HashMap<>();
				fileMsgMap.put("role", "system");
				fileMsgMap.put("content", (String)fileContent.get("content"));
				messageList.add(fileMsgMap);
			}
		}
		if(StringUtils.isNotBlank(goOnContent) && StringUtils.isNotBlank(assistantMap.get(goOnContent))){
			Map<String, String> filegoOnMsgMap = new HashMap<>();
			filegoOnMsgMap.put("role", "system");
			filegoOnMsgMap.put("content", assistantMap.get(goOnContent));
			messageList.add(filegoOnMsgMap);
		}
//		assistantMap.put("role", "assistant");
//		assistantMap.put("content", "");
//		StringBuilder temp = new StringBuilder();
		String result = null;
		// 发起异步请求
		try {
			result = HttpClientUtil.postAddHeadersWithObjEntity(requestUrl, paramMap, headerMap);
			//从6开始 因为有 data: 这个前缀 占了6个字符所以 0 + 6 = 6
			jsonObject = JSON.parseObject(result);
			if (jsonObject != null && jsonObject.getJSONArray("choices") != null && jsonObject.getJSONArray("choices").size() > 0) {
				result = jsonObject.getJSONArray("choices").getJSONObject(0).getJSONObject("message").getString("content");
			}
			if (StringUtils.isNotBlank(result)) {
				//SSE协议默认是以两个\n换行符为结束标志 需要在进行一次转义才能成功发送给前端
				//将结果汇总起来
				assistantMap.put(clientId, result);
			}
		} catch (Exception e) {
			System.out.println("___--------" + result);
			e.printStackTrace();
		}
		System.out.println("___--------" + result);
		System.out.println("___--------" + jsonObject.toJSONString());
		return jsonObject;
	}

}
