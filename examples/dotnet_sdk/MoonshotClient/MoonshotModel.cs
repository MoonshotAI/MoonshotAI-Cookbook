using System.Collections.Generic;

namespace MoonshotDotnet
{


    public class MessagesItem
    {
        /// <summary>
        /// 
        /// </summary>
        public string role { get; set; }

        /// <summary>
        /// 
        /// </summary>
        public string content { get; set; }
    }

    public class ChatReq
    {
        /// <summary>
        /// 
        /// </summary>
        public string model { get; set; }
        /// <summary>
        /// 
        /// </summary>
        public List<MessagesItem> messages { get; set; }
        /// <summary>
        /// 
        /// </summary>
        public double temperature { get; set; }


        public int? max_tokens { get; set; }

        public float? top_p { get; set; } = 1.0f;

        public int? n { get; set; } = 1;

        public float? presence_penalty { get; set; } = 0;

        public float? frequency_penalty { get; set; } = 0;

        public List<string> stop { get; set; }

        public bool stream { get; set; } = false;


    }


    public class PermissionItem
    {
        /// <summary>
        /// 
        /// </summary>
        public int created { get; set; }
        /// <summary>
        /// 
        /// </summary>
        public string id { get; set; }
        /// <summary>
        /// 
        /// </summary>
        public string @object { get; set; }
        /// <summary>
        /// 
        /// </summary>
        public string allow_create_engine { get; set; }
        /// <summary>
        /// 
        /// </summary>
        public string allow_sampling { get; set; }
        /// <summary>
        /// 
        /// </summary>
        public string allow_logprobs { get; set; }
        /// <summary>
        /// 
        /// </summary>
        public string allow_search_indices { get; set; }
        /// <summary>
        /// 
        /// </summary>
        public string allow_view { get; set; }
        /// <summary>
        /// 
        /// </summary>
        public string allow_fine_tuning { get; set; }
        /// <summary>
        /// 
        /// </summary>
        public string organization { get; set; }
        /// <summary>
        /// 
        /// </summary>
        public string @group { get; set; }
        /// <summary>
        /// 
        /// </summary>
        public string is_blocking { get; set; }
    }

    public class ModelInfo
    {
        /// <summary>
        /// 
        /// </summary>
        public int created { get; set; }
        /// <summary>
        /// 
        /// </summary>
        public string id { get; set; }
        /// <summary>
        /// 
        /// </summary>
        public string @object { get; set; }
        /// <summary>
        /// 
        /// </summary>
        public string owned_by { get; set; }
        /// <summary>
        /// 
        /// </summary>
        public List<PermissionItem> permission { get; set; }
        /// <summary>
        /// 
        /// </summary>
        public string root { get; set; }
        /// <summary>
        /// 
        /// </summary>
        public string parent { get; set; }
    }

    public class ModelListResp
    {
        /// <summary>
        /// 
        /// </summary>
        public string @object { get; set; }
        /// <summary>
        /// 
        /// </summary>
        public List<ModelInfo> data { get; set; }
    }


    public class FileListResp
    {
        /// <summary>
        /// 
        /// </summary>
        public string @object { get; set; }
        /// <summary>
        /// 
        /// </summary>
        public List<FileItem> data { get; set; }
    }

    public class FileContent
    {
        /// <summary>
        /// 
        /// </summary>
        public string content { get; set; }
        /// <summary>
        /// 
        /// </summary>
        public string file_type { get; set; }
        /// <summary>
        /// 
        /// </summary>
        public string filename { get; set; }
        /// <summary>
        /// 
        /// </summary>
        public string title { get; set; }
        /// <summary>
        /// 
        /// </summary>
        public string type { get; set; }
    }


    public class FileItem
    {
        /// <summary>
        /// 
        /// </summary>
        public string id { get; set; }
        /// <summary>
        /// 
        /// </summary>
        public string @object { get; set; }
        /// <summary>
        /// 
        /// </summary>
        public int bytes { get; set; }
        /// <summary>
        /// 
        /// </summary>
        public int created_at { get; set; }
        /// <summary>
        /// 
        /// </summary>
        public string filename { get; set; }
        /// <summary>
        /// 
        /// </summary>
        public string purpose { get; set; }
        /// <summary>
        /// 
        /// </summary>
        public string status { get; set; }
        /// <summary>
        /// 
        /// </summary>
        public string status_details { get; set; }
    }



    public class ErrorMsg
    {
        /// <summary>
        /// 
        /// </summary>
        public string message { get; set; }
        /// <summary>
        /// 
        /// </summary>
        public string type { get; set; }
    }

    public class ErrorResponse
    {
        /// <summary>
        /// 
        /// </summary>
        public ErrorMsg error { get; set; }
    }



}
