package cn.moonshot.platform.util;

import cn.hutool.core.lang.Assert;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.EqualsAndHashCode;
import lombok.NoArgsConstructor;

@NoArgsConstructor
@AllArgsConstructor
@Data
@Builder
public class Message {

    private String role;

    private String content;

}
