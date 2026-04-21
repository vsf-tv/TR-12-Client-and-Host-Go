import { ResponseContext, RequestContext, HttpFile, HttpInfo } from '../http/http';
import { Configuration, ConfigurationOptions, mergeConfiguration } from '../configuration'
import type { Middleware } from '../middleware';
import { Observable, of, from } from '../rxjsStub';
import {mergeMap, map} from  '../rxjsStub';
import { Aes128 } from '../models/Aes128';
import { Aes256 } from '../models/Aes256';
import { Channel } from '../models/Channel';
import { ChannelConfiguration } from '../models/ChannelConfiguration';
import { ChannelProfile } from '../models/ChannelProfile';
import { ChannelState } from '../models/ChannelState';
import { ChannelStatus } from '../models/ChannelStatus';
import { ChannelType } from '../models/ChannelType';
import { ConfigurationData } from '../models/ConfigurationData';
import { ConnectRequestContent } from '../models/ConnectRequestContent';
import { ConnectResponseContent } from '../models/ConnectResponseContent';
import { Connection } from '../models/Connection';
import { Critical } from '../models/Critical';
import { Degraded } from '../models/Degraded';
import { DeprovisionResponseContent } from '../models/DeprovisionResponseContent';
import { DeviceConfiguration } from '../models/DeviceConfiguration';
import { DeviceRegistration } from '../models/DeviceRegistration';
import { DeviceStatus } from '../models/DeviceStatus';
import { DisconnectResponseContent } from '../models/DisconnectResponseContent';
import { EncryptionAes } from '../models/EncryptionAes';
import { EncryptionAes128 } from '../models/EncryptionAes128';
import { EncryptionAes256 } from '../models/EncryptionAes256';
import { EnumValues } from '../models/EnumValues';
import { ErrorDetails } from '../models/ErrorDetails';
import { GetConfigurationResponseContent } from '../models/GetConfigurationResponseContent';
import { GetConnectionStatusResponseContent } from '../models/GetConnectionStatusResponseContent';
import { Health } from '../models/Health';
import { Healthy } from '../models/Healthy';
import { IdAndValue } from '../models/IdAndValue';
import { Profile } from '../models/Profile';
import { ProfileDefinition } from '../models/ProfileDefinition';
import { RangeValues } from '../models/RangeValues';
import { ReportActualConfigurationRequestContent } from '../models/ReportActualConfigurationRequestContent';
import { ReportActualConfigurationResponseContent } from '../models/ReportActualConfigurationResponseContent';
import { ReportStatusRequestContent } from '../models/ReportStatusRequestContent';
import { ReportStatusResponseContent } from '../models/ReportStatusResponseContent';
import { RistCaller } from '../models/RistCaller';
import { RistCallerTransportProtocol } from '../models/RistCallerTransportProtocol';
import { RistListener } from '../models/RistListener';
import { RistListenerTransportProtocol } from '../models/RistListenerTransportProtocol';
import { RistStreamIdentifier } from '../models/RistStreamIdentifier';
import { Rtp } from '../models/Rtp';
import { RtpFecConfiguration } from '../models/RtpFecConfiguration';
import { RtpFecStreamConfig } from '../models/RtpFecStreamConfig';
import { RtpTransportProtocol } from '../models/RtpTransportProtocol';
import { Setting } from '../models/Setting';
import { SettingsChoice } from '../models/SettingsChoice';
import { SrtCaller } from '../models/SrtCaller';
import { SrtCallerTransportProtocol } from '../models/SrtCallerTransportProtocol';
import { SrtListener } from '../models/SrtListener';
import { SrtListenerTransportProtocol } from '../models/SrtListenerTransportProtocol';
import { StandardSettings } from '../models/StandardSettings';
import { StatusValue } from '../models/StatusValue';
import { StreamId } from '../models/StreamId';
import { SynchronizationSource } from '../models/SynchronizationSource';
import { Thumbnail } from '../models/Thumbnail';
import { TransportProtocol } from '../models/TransportProtocol';
import { TransportProtocolName } from '../models/TransportProtocolName';
import { UnhealthyStateDescription } from '../models/UnhealthyStateDescription';
import { ZixiPull } from '../models/ZixiPull';
import { ZixiPullTransportProtocol } from '../models/ZixiPullTransportProtocol';
import { ZixiPush } from '../models/ZixiPush';
import { ZixiPushTransportProtocol } from '../models/ZixiPushTransportProtocol';

import { DefaultApiRequestFactory, DefaultApiResponseProcessor} from "../apis/DefaultApi";
export class ObservableDefaultApi {
    private requestFactory: DefaultApiRequestFactory;
    private responseProcessor: DefaultApiResponseProcessor;
    private configuration: Configuration;

    public constructor(
        configuration: Configuration,
        requestFactory?: DefaultApiRequestFactory,
        responseProcessor?: DefaultApiResponseProcessor
    ) {
        this.configuration = configuration;
        this.requestFactory = requestFactory || new DefaultApiRequestFactory(configuration);
        this.responseProcessor = responseProcessor || new DefaultApiResponseProcessor();
    }

    /**
     * @param connectRequestContent
     */
    public connectWithHttpInfo(connectRequestContent: ConnectRequestContent, _options?: ConfigurationOptions): Observable<HttpInfo<ConnectResponseContent>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.connect(connectRequestContent, _config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.connectWithHttpInfo(rsp)));
            }));
    }

    /**
     * @param connectRequestContent
     */
    public connect(connectRequestContent: ConnectRequestContent, _options?: ConfigurationOptions): Observable<ConnectResponseContent> {
        return this.connectWithHttpInfo(connectRequestContent, _options).pipe(map((apiResponse: HttpInfo<ConnectResponseContent>) => apiResponse.data));
    }

    /**
     * @param hostId
     * @param [force]
     */
    public deprovisionWithHttpInfo(hostId: string, force?: boolean, _options?: ConfigurationOptions): Observable<HttpInfo<DeprovisionResponseContent>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.deprovision(hostId, force, _config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.deprovisionWithHttpInfo(rsp)));
            }));
    }

    /**
     * @param hostId
     * @param [force]
     */
    public deprovision(hostId: string, force?: boolean, _options?: ConfigurationOptions): Observable<DeprovisionResponseContent> {
        return this.deprovisionWithHttpInfo(hostId, force, _options).pipe(map((apiResponse: HttpInfo<DeprovisionResponseContent>) => apiResponse.data));
    }

    /**
     */
    public disconnectWithHttpInfo(_options?: ConfigurationOptions): Observable<HttpInfo<DisconnectResponseContent>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.disconnect(_config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.disconnectWithHttpInfo(rsp)));
            }));
    }

    /**
     */
    public disconnect(_options?: ConfigurationOptions): Observable<DisconnectResponseContent> {
        return this.disconnectWithHttpInfo(_options).pipe(map((apiResponse: HttpInfo<DisconnectResponseContent>) => apiResponse.data));
    }

    /**
     */
    public getConfigurationWithHttpInfo(_options?: ConfigurationOptions): Observable<HttpInfo<GetConfigurationResponseContent>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.getConfiguration(_config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.getConfigurationWithHttpInfo(rsp)));
            }));
    }

    /**
     */
    public getConfiguration(_options?: ConfigurationOptions): Observable<GetConfigurationResponseContent> {
        return this.getConfigurationWithHttpInfo(_options).pipe(map((apiResponse: HttpInfo<GetConfigurationResponseContent>) => apiResponse.data));
    }

    /**
     */
    public getConnectionStatusWithHttpInfo(_options?: ConfigurationOptions): Observable<HttpInfo<GetConnectionStatusResponseContent>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.getConnectionStatus(_config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.getConnectionStatusWithHttpInfo(rsp)));
            }));
    }

    /**
     */
    public getConnectionStatus(_options?: ConfigurationOptions): Observable<GetConnectionStatusResponseContent> {
        return this.getConnectionStatusWithHttpInfo(_options).pipe(map((apiResponse: HttpInfo<GetConnectionStatusResponseContent>) => apiResponse.data));
    }

    /**
     * @param reportActualConfigurationRequestContent
     */
    public reportActualConfigurationWithHttpInfo(reportActualConfigurationRequestContent: ReportActualConfigurationRequestContent, _options?: ConfigurationOptions): Observable<HttpInfo<ReportActualConfigurationResponseContent>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.reportActualConfiguration(reportActualConfigurationRequestContent, _config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.reportActualConfigurationWithHttpInfo(rsp)));
            }));
    }

    /**
     * @param reportActualConfigurationRequestContent
     */
    public reportActualConfiguration(reportActualConfigurationRequestContent: ReportActualConfigurationRequestContent, _options?: ConfigurationOptions): Observable<ReportActualConfigurationResponseContent> {
        return this.reportActualConfigurationWithHttpInfo(reportActualConfigurationRequestContent, _options).pipe(map((apiResponse: HttpInfo<ReportActualConfigurationResponseContent>) => apiResponse.data));
    }

    /**
     * @param reportStatusRequestContent
     */
    public reportStatusWithHttpInfo(reportStatusRequestContent: ReportStatusRequestContent, _options?: ConfigurationOptions): Observable<HttpInfo<ReportStatusResponseContent>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.reportStatus(reportStatusRequestContent, _config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.reportStatusWithHttpInfo(rsp)));
            }));
    }

    /**
     * @param reportStatusRequestContent
     */
    public reportStatus(reportStatusRequestContent: ReportStatusRequestContent, _options?: ConfigurationOptions): Observable<ReportStatusResponseContent> {
        return this.reportStatusWithHttpInfo(reportStatusRequestContent, _options).pipe(map((apiResponse: HttpInfo<ReportStatusResponseContent>) => apiResponse.data));
    }

}
